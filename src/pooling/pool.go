package pool

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const _Beginner = "beginner"
const _Novice = "novice"
const _LowerIntermediate = "lower_intermediate"
const _Intermediate = "intermediate"
const _UpperIntermediate = "intermediate"
const _Advanced = "advanced"
const _Expert = "expert"
const _Godlike = "godlike"

// Request encapsulates pool data.
type Request struct {
	ConnectionID string
	UserID       string
	PoolLimit    int
}

// JoinOrCreate joins a user to an existing pool
// (relative to their skill level), or creates a new one.
func JoinOrCreate(r *Request) {
	dynamo := dynamo()
	bucket := getPoolBucket(&r.UserID)
	poolIDPtr := getAvailablePoolID(dynamo, bucket)
	var poolID string

	if poolIDPtr == nil {
		poolID = new(bucket, r, dynamo)

		addPoolToBucket(bucket, poolID, dynamo)
	} else {
		poolID = *poolIDPtr

		joinPool(poolID, r.ConnectionID, dynamo)
	}

	mapConnectionToPool(r.ConnectionID, poolID, dynamo)
}

// Leave removes a connection from a given pool.
func Leave(connectionID string) error {
	dynamo := dynamo()
	poolID, err := getPoolID(connectionID, dynamo)

	if err != nil {
		return err
	}

	key := map[string]*dynamodb.AttributeValue{
		"ID": {
			S: aws.String(poolID),
		},
	}
	ue := aws.String("DELETE ConnectionIDs :cids")
	eav := map[string]*dynamodb.AttributeValue{
		":cids": {
			SS: []*string{aws.String(connectionID)},
		},
	}
	_, updateErr := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		Key:                       key,
		UpdateExpression:          ue,
		ExpressionAttributeValues: eav,
		TableName:                 aws.String("pools"),
	})

	if updateErr != nil {
		return updateErr
	}

	_, deleteErr := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(connectionID),
			},
		},
		TableName: aws.String("connections"),
	})

	return deleteErr
}

func dynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-central-1")}))
	dynamo := dynamodb.New(sess)
	return dynamo
}

func getPoolBucket(userID *string) string {
	if userID == nil {
		return _Beginner
	}

	// Look up user's level
	return _Novice
}

type bucketItem struct {
	CurrentAvailablePool string
}

func getAvailablePoolID(dynamo *dynamodb.DynamoDB, bucket string) *string {
	o, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(bucket),
			},
		},
		ConsistentRead:  aws.Bool(true),
		AttributesToGet: []*string{aws.String("CurrentAvailablePool")},
		TableName:       aws.String("buckets"),
	})

	if err != nil || o.Item == nil {
		fmt.Println(err)
		return nil
	}

	bucketItem := bucketItem{}
	err = dynamodbattribute.UnmarshalMap(o.Item, &bucketItem)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if !poolHasCapacity(bucketItem.CurrentAvailablePool, dynamo) {
		return nil
	}

	return &bucketItem.CurrentAvailablePool
}

func new(bucket string, r *Request, dynamo *dynamodb.DynamoDB) string {
	poolID := uuid.New().String()
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(poolID),
			},
			"ConnectionIDs": {
				SS: []*string{aws.String(r.ConnectionID)},
			},
			"Limit": {
				N: aws.String(strconv.Itoa(r.PoolLimit)),
			},
			"Bucket": {
				S: aws.String(bucket),
			},
		},
		TableName: aws.String("pools"),
	})

	if err != nil {
		fmt.Println(err)
	}

	return poolID
}

type poolItem struct {
	ID            string
	ConnectionIDs []string
	Limit         int
}

func poolHasCapacity(poolID string, dynamo *dynamodb.DynamoDB) bool {
	o, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(poolID),
			},
		},
		ConsistentRead: aws.Bool(true),
		TableName:      aws.String("pools"),
	})

	if err != nil || o.Item == nil {
		fmt.Println(err)
		return false
	}

	item := poolItem{}
	dynamodbattribute.UnmarshalMap(o.Item, &item)

	return len(item.ConnectionIDs) < item.Limit
}

func addPoolToBucket(bucket string, poolID string, dynamo *dynamodb.DynamoDB) {
	key := map[string]*dynamodb.AttributeValue{
		"ID": {
			S: aws.String(bucket),
		},
	}
	ue := aws.String("ADD PoolIDs :pids SET CurrentAvailablePool = :cap")
	eav := map[string]*dynamodb.AttributeValue{
		":pids": {
			SS: []*string{aws.String(poolID)},
		},
		":cap": {
			S: aws.String(poolID),
		},
	}
	_, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		Key:                       key,
		UpdateExpression:          ue,
		ExpressionAttributeValues: eav,
		TableName:                 aws.String("buckets"),
	})

	if err != nil {
		fmt.Println(err)
		// TODO Do something...
	}
}

func joinPool(poolID string, connectionID string, dynamo *dynamodb.DynamoDB) {
	key := map[string]*dynamodb.AttributeValue{
		"ID": {
			S: aws.String(poolID),
		},
	}
	ue := aws.String("ADD ConnectionIDs :cids")
	eav := map[string]*dynamodb.AttributeValue{
		":cids": {
			SS: []*string{aws.String(connectionID)},
		},
	}
	_, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		Key:                       key,
		UpdateExpression:          ue,
		ExpressionAttributeValues: eav,
		TableName:                 aws.String("pools"),
	})

	if err != nil {
		fmt.Println(err)
		// TODO Do something...
	}
}

type connection struct {
	PoolID string
}

func getPoolID(connectionID string, dynamo *dynamodb.DynamoDB) (string, error) {
	o, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(connectionID),
			},
		},
		ConsistentRead:  aws.Bool(true),
		AttributesToGet: []*string{aws.String("PoolID")},
		TableName:       aws.String("connections"),
	})

	if err != nil {
		return "", err
	}

	if o.Item == nil {
		return "", nil
	}

	conn := connection{}
	err = dynamodbattribute.UnmarshalMap(o.Item, &conn)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return conn.PoolID, nil
}

func mapConnectionToPool(connectionID, poolID string, dynamo *dynamodb.DynamoDB) error {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(connectionID),
			},
			"PoolID": {
				S: aws.String(poolID),
			},
		},
		TableName: aws.String("connections"),
	})

	return err
}
