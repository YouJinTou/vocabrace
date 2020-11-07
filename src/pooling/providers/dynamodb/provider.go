package dynamodbpooling

import (
	"errors"
	"fmt"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBProvider implements Provider.
type DynamoDBProvider struct {
	dynamo *dynamodb.DynamoDB
	stage  string
}

type connection struct {
	ID     string
	PoolID string
}

// NewDynamoDBProvider creates a new pooling provider using DynamoDB as a backend.
func NewDynamoDBProvider(stage string) pooling.Provider {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)

	return DynamoDBProvider{
		dynamo: dynamo,
		stage:  stage,
	}
}

// GetPool gets a pool.
func (dpp DynamoDBProvider) GetPool(i *pooling.GetPoolInput) (*pooling.Pool, error) {
	result, err := dpp.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_buckets", dpp.stage)),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(i.Bucket)},
		},
		ExpressionAttributeNames: map[string]*string{"#p": aws.String(i.PoolID)},
		ProjectionExpression:     aws.String("#p"),
	})
	p := pooling.Pool{Bucket: i.Bucket}
	dynamodbattribute.UnmarshalMap(result.Item[i.PoolID].M, &p)

	if p.ID == "" {
		return nil, errors.New("could not find pool")
	}

	return &p, err
}

// GetPeers gets a connection's pool peers.
func (dpp DynamoDBProvider) GetPeers(i *pooling.GetPeersInput) ([]string, error) {
	c, cErr := dpp.getConnection(i.ConnectionID)
	peers := []string{}

	if cErr != nil {
		return peers, cErr
	}

	p, pErr := dpp.GetPool(&pooling.GetPoolInput{
		Bucket: i.Bucket,
		PoolID: c.PoolID,
	})

	if pErr != nil {
		return peers, pErr
	}

	for _, cid := range p.ConnectionIDs {
		if i.ConnectionID != cid {
			peers = append(peers, cid)
		}
	}

	return peers, nil
}

func (dpp DynamoDBProvider) getConnection(connectionID string) (connection, error) {
	result, err := dpp.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_connections", dpp.stage)),
		Key:       map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(connectionID)}},
	})
	c := connection{}
	dynamodbattribute.UnmarshalMap(result.Item, &c)

	if c.ID == "" {
		return c, errors.New("could not find connection")
	}

	return c, err
}
