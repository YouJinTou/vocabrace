package pool

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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
