package ws

import (
	"fmt"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// GetPoolID gets a pool ID given a connection ID.
func GetPoolID(connectionID string, c *Config) (*string, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	o, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName:            aws.String(fmt.Sprintf("%s_connections", c.Stage)),
		Key:                  map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(connectionID)}},
		ProjectionExpression: aws.String("PoolID"),
	})

	return o.Item["PoolID"].S, err
}

// GetPool using a connectionID
func GetPool(connectionID string, c *Config) (*pooling.Pool, error) {
	poolID, gErr := GetPoolID(connectionID, c)

	if poolID == nil {
		return nil, gErr
	}

	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	i, pErr := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_pools", c.Stage)),
		Key:       map[string]*dynamodb.AttributeValue{"ID": {S: poolID}},
	})
	pool := pooling.Pool{}
	dynamodbattribute.UnmarshalMap(i.Item, &pool)

	return &pool, pErr
}
