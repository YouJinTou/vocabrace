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

// Request encapsulates pool data
type Request struct {
	ConnectionID string
	UserID       string
	PoolLimit    int
}

// JoinOrCreate joins a user to an existing pool
// (relative to their skill level), or creates a new one
func JoinOrCreate(r *Request) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-central-1")}))
	dynamo := dynamodb.New(sess)
	bucket := getPoolBucket(&r.UserID)
	var poolIDPtr = getAvailablePoolID(dynamo, bucket)
	var poolID string

	if poolIDPtr == nil {
		poolID = createNewPool(bucket, r, dynamo)

		addPoolToBucket(bucket, poolID, dynamo)
	} else {
		poolID = *poolIDPtr

		if !poolHasCapacity(poolID, dynamo) {
			poolID = createNewPool(bucket, r, dynamo)
		}
	}

	joinPool(poolID, r.ConnectionID, dynamo)
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

	return &bucketItem.CurrentAvailablePool
}

func createNewPool(bucket string, r *Request, dynamo *dynamodb.DynamoDB) string {
	poolID := uuid.New().String()
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(poolID),
			},
			"ConnectionIDs": {
				L: []*dynamodb.AttributeValue{},
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
		TableName: aws.String("pools"),
	})

	if err != nil || o.Item == nil {
		fmt.Println(err)
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
	ue := aws.String("SET PoolIDs = list_append(if_not_exists(PoolIDs, :empty_list), :pids), CurrentAvailablePool = :cap")
	eav := map[string]*dynamodb.AttributeValue{
		":pids": {
			L: []*dynamodb.AttributeValue{
				&dynamodb.AttributeValue{S: aws.String(poolID)},
			},
		},
		":empty_list": {
			L: []*dynamodb.AttributeValue{},
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
	ue := aws.String("SET ConnectionIDs = list_append(ConnectionIDs, :cids)")
	eav := map[string]*dynamodb.AttributeValue{
		":cids": {
			L: []*dynamodb.AttributeValue{
				&dynamodb.AttributeValue{S: aws.String(connectionID)},
			},
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
