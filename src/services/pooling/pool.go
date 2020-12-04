package pooling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Pool carries pool data.
type Pool struct {
	ID            string
	ConnectionIDs []string
	Bucket        string
	Limit         int
}

// Beginner is the beginner bucket.
const Beginner = "beginner"

// Novice is the novice bucket.
const Novice = "novice"

// LowerIntermediate is the lower_intermediate bucket.
const LowerIntermediate = "lower_intermediate"

// Intermediate is the intermediate bucket.
const Intermediate = "intermediate"

// UpperIntermediate is the upper_intermediate bucket.
const UpperIntermediate = "upper_intermediate"

// Advanced is the advanced bucket.
const Advanced = "advanced"

// Expert is the expert bucket.
const Expert = "expert"

// Godlike is the godlike bucket.
const Godlike = "godlike"

// GetPoolID gets a pool ID given a connection ID.
func GetPoolID(connectionID, stage string) (*string, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	o, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName:            aws.String(fmt.Sprintf("%s_connections", stage)),
		Key:                  map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(connectionID)}},
		ProjectionExpression: aws.String("PoolID"),
	})

	return o.Item["PoolID"].S, err
}

// GetPool using a connectionID
func GetPool(connectionID, stage string) (*Pool, error) {
	poolID, gErr := GetPoolID(connectionID, stage)

	if poolID == nil {
		return nil, gErr
	}

	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	i, pErr := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_pools", stage)),
		Key:       map[string]*dynamodb.AttributeValue{"ID": {S: poolID}},
	})
	pool := Pool{}
	dynamodbattribute.UnmarshalMap(i.Item, &pool)

	return &pool, pErr
}
