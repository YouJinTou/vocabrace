package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/YouJinTou/vocabrace/tools"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Pool carries pool data.
type Pool struct {
	ID            string
	ConnectionIDs []string
	Bucket        string
	Limit         int
}

type dbBucket struct {
	ID        string
	CAP       *string
	UpdatedAt int
}

type dbBucketWrapper struct {
	dbb        *dbBucket
	CAPPool    *dbPool
	CAPCreated bool
}

type dbPool struct {
	ID            string
	ConnectionIDs []string
	Limit         int
}

func (dbp *dbBucket) exists() bool {
	return dbp.ID != ""
}

func (dbp *dbPool) isFull() bool {
	return len(dbp.ConnectionIDs) >= dbp.Limit
}

// Request encapsulates pool data.
type Request struct {
	ConnectionID string
	UserID       string
	Bucket       string
	PoolLimit    int
	Stage        string
}

// PoolingProvider abstracts a pooling provider.
type PoolingProvider interface {
	JoinOrCreate(r *Request) (*Pool, error)
}

// DynamoDBPoolingProvider implements PoolingProvider.
type DynamoDBPoolingProvider struct {
}

// NewDynamoDBPoolingProvider creates a new DynamoDBPoolingProvider.
func NewDynamoDBPoolingProvider() PoolingProvider {
	return DynamoDBPoolingProvider{}
}

func main() {
	for i := 0; i < 11; i++ {
		go NewDynamoDBPoolingProvider().JoinOrCreate(&Request{
			ConnectionID: uuid.New().String(),
			UserID:       uuid.New().String(),
			Bucket:       "novice",
			PoolLimit:    3,
			Stage:        "dev",
		})
	}
	time.Sleep(10 * time.Second)
}

// JoinOrCreate joins or creates a pool.
func (dpp DynamoDBPoolingProvider) JoinOrCreate(r *Request) (*Pool, error) {
	for {
		w, getErr := getDbBucketWrapper(r)

		if getErr != nil {
			continue
		}

		if !w.dbb.exists() {
			newB, err := createDbBucket(r)
			if err == nil {
				w.dbb = newB
			}
		}

		mapConnection(w, r)

		p, setErr := setPool(w, r)

		if setErr != nil {
			continue
		}

		return p, nil
	}
}

func getDbBucketWrapper(r *Request) (*dbBucketWrapper, error) {
	result, err := dynamo().GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_buckets", r.Stage)),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(r.Bucket),
			},
		},
		ProjectionExpression: aws.String("ID, CAP, UpdatedAt"),
	})
	w := dbBucketWrapper{}
	dynamodbattribute.UnmarshalMap(result.Item, &w.dbb)

	if w.dbb.CAP != nil {
		result, _ := dynamo().GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(fmt.Sprintf("%s_buckets", r.Stage)),
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(r.Bucket),
				},
			},
			ExpressionAttributeNames: map[string]*string{"#cap": w.dbb.CAP},
			ProjectionExpression:     aws.String("#cap"),
		})
		cap := dbPool{}
		av := result.Item[*w.dbb.CAP]
		dynamodbattribute.UnmarshalMap(av.M, &cap)
		w.CAPPool = &cap
	}

	return &w, err
}

func createDbBucket(r *Request) (*dbBucket, error) {
	b := dbBucket{
		ID:        r.Bucket,
		CAP:       nil,
		UpdatedAt: tools.FutureTimestamp(0),
	}
	marshaled, _ := dynamodbattribute.MarshalMap(b)
	_, err := dynamo().PutItem(&dynamodb.PutItemInput{
		TableName:           aws.String(fmt.Sprintf("%s_buckets", r.Stage)),
		Item:                marshaled,
		ConditionExpression: aws.String("attribute_not_exists(UpdatedAt)"),
	})

	return &b, err
}

func mapConnection(w *dbBucketWrapper, r *Request) {
	cap := w.CAPPool

	if cap == nil || cap.isFull() {
		ID := fmt.Sprintf("ZZ%s", uuid.New().String()[0:5])
		w.CAPPool = &dbPool{
			ID:            ID,
			ConnectionIDs: []string{r.ConnectionID},
			Limit:         r.PoolLimit,
		}
		w.dbb.CAP = &ID
		w.CAPCreated = true
	} else {
		cap.ConnectionIDs = append(cap.ConnectionIDs, r.ConnectionID)
	}
}

func setPool(w *dbBucketWrapper, r *Request) (*Pool, error) {
	key := map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(r.Bucket)}}
	eav := map[string]*dynamodb.AttributeValue{
		":ua":  {N: aws.String(tools.FutureTimestampStr(86400))},
		":lua": {N: aws.String(strconv.Itoa(w.dbb.UpdatedAt))},
	}
	var ue string

	if w.CAPCreated {
		ue = "SET CAP = :cap, UpdatedAt = :ua, #p = :p"
		eav[":cap"] = &dynamodb.AttributeValue{S: w.dbb.CAP}
		eav[":p"] = &dynamodb.AttributeValue{
			M: map[string]*dynamodb.AttributeValue{
				"ID":            {S: w.dbb.CAP},
				"ConnectionIDs": {SS: []*string{aws.String(r.ConnectionID)}},
				"Limit":         {N: aws.String(strconv.Itoa(r.PoolLimit))},
			},
		}
	} else {
		ue = "SET UpdatedAt = :ua ADD #p.ConnectionIDs :cid"
		eav[":cid"] = &dynamodb.AttributeValue{SS: []*string{aws.String(r.ConnectionID)}}
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(fmt.Sprintf("%s_buckets", r.Stage)),
		Key:                       key,
		ExpressionAttributeValues: eav,
		ExpressionAttributeNames:  map[string]*string{"#p": w.dbb.CAP},
		UpdateExpression:          aws.String(ue),
		ConditionExpression:       aws.String("UpdatedAt = :lua"),
		ReturnValues:              aws.String("ALL_NEW"),
	}
	result, err := dynamo().UpdateItem(input)
	pool := Pool{}
	dynamodbattribute.UnmarshalMap(result.Attributes[*w.dbb.CAP].M, &pool)
	pool.Bucket = r.Bucket

	return &pool, err
}

func dynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}
