package dynamodbpooling

import (
	"fmt"
	"strconv"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/YouJinTou/vocabrace/tools"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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

// JoinOrCreate joins or creates a pool.
func (dpp DynamoDBProvider) JoinOrCreate(r *pooling.Request) (*pooling.Pool, error) {
	for {
		w, getErr := dpp.getDbBucketWrapper(r)

		if getErr != nil {
			continue
		}

		if !w.dbb.exists() {
			newB, err := dpp.createDbBucket(r)
			if err == nil {
				w.dbb = newB
			}
		}

		dpp.mapConnection(w, r)

		p, setErr := dpp.setPool(w, r)

		if setErr != nil {
			continue
		}

		dpp.setConnection(p.ID, r)

		return p, nil
	}
}

func (dpp DynamoDBProvider) getDbBucket(r *pooling.Request) (*dbBucket, error) {
	result, err := dpp.dynamo().GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_buckets", r.Stage)),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(r.Bucket),
			},
		},
		ProjectionExpression: aws.String("ID, CAP, UpdatedAt"),
	})
	b := dbBucket{}
	dynamodbattribute.UnmarshalMap(result.Item, &b)

	return &b, err
}

func (dpp DynamoDBProvider) getDbBucketWrapper(r *pooling.Request) (*dbBucketWrapper, error) {
	dbb, err := dpp.getDbBucket(r)
	w := dbBucketWrapper{dbb: dbb}

	if w.dbb.CAP != nil {
		result, _ := dpp.dynamo().GetItem(&dynamodb.GetItemInput{
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

func (dpp DynamoDBProvider) createDbBucket(r *pooling.Request) (*dbBucket, error) {
	b := dbBucket{
		ID:        r.Bucket,
		CAP:       nil,
		UpdatedAt: tools.FutureTimestamp(0),
	}
	marshaled, _ := dynamodbattribute.MarshalMap(b)
	_, err := dpp.dynamo().PutItem(&dynamodb.PutItemInput{
		TableName:           aws.String(fmt.Sprintf("%s_buckets", r.Stage)),
		Item:                marshaled,
		ConditionExpression: aws.String("attribute_not_exists(UpdatedAt)"),
	})

	return &b, err
}

func (dpp DynamoDBProvider) mapConnection(w *dbBucketWrapper, r *pooling.Request) {
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

func (dpp DynamoDBProvider) setPool(w *dbBucketWrapper, r *pooling.Request) (*pooling.Pool, error) {
	key := map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(r.Bucket)}}
	eav := map[string]*dynamodb.AttributeValue{
		":ua":  {N: aws.String(tools.FutureTimestampStr(0))},
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
	result, err := dpp.dynamo().UpdateItem(input)
	pool := pooling.Pool{
		Bucket: r.Bucket,
	}
	dynamodbattribute.UnmarshalMap(result.Attributes[*w.dbb.CAP].M, &pool)

	return &pool, err
}

func (dpp DynamoDBProvider) setConnection(poolID string, r *pooling.Request) {
	dpp.dynamo().UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(fmt.Sprintf("%s_connections", r.Stage)),
		Key:       map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(r.ConnectionID)}},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pid": {S: aws.String(poolID)},
		},
		UpdateExpression: aws.String("SET PoolID = :pid"),
	})
}
