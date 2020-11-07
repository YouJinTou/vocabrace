package dynamodbpooling

import (
	"fmt"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBProvider implements Provider.
type DynamoDBProvider struct {
	d *dynamodb.DynamoDB
}

type connection struct {
	ID     string
	PoolID string
}

// NewDynamoDBProvider creates a new pooling provider using DynamoDB as a backend.
func NewDynamoDBProvider() pooling.Provider {
	return DynamoDBProvider{}
}

// GetPool gets a pool.
func (dpp DynamoDBProvider) GetPool(ID string, r *pooling.Request) (*pooling.Pool, error) {
	result, err := dpp.dynamo().GetItem(&dynamodb.GetItemInput{
		TableName:                aws.String(fmt.Sprintf("%s_buckets", r.Stage)),
		Key:                      map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(r.Bucket)}},
		ExpressionAttributeNames: map[string]*string{"#p": aws.String(ID)},
		ProjectionExpression:     aws.String("#p"),
	})
	p := pooling.Pool{}
	dynamodbattribute.UnmarshalMap(result.Item, &p)

	return &p, err
}

// GetPeers gets a connection's pool peers.
func (dpp DynamoDBProvider) GetPeers(r *pooling.Request) ([]string, error) {
	c, cErr := dpp.getConnection(r)
	peers := []string{}

	if cErr != nil {
		return peers, cErr
	}

	p, pErr := dpp.GetPool(c.PoolID, r)

	if pErr != nil {
		return peers, pErr
	}

	for _, cid := range p.ConnectionIDs {
		if r.ConnectionID != cid {
			peers = append(peers, cid)
		}
	}

	return peers, nil
}

func (dpp DynamoDBProvider) getConnection(r *pooling.Request) (connection, error) {
	result, err := dpp.dynamo().GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_connections", r.Stage)),
		Key:       map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(r.ConnectionID)}},
	})
	c := connection{}
	dynamodbattribute.UnmarshalMap(result.Item, &c)

	return c, err
}

func (dpp DynamoDBProvider) dynamo() *dynamodb.DynamoDB {
	if dpp.d == nil {
		sess := session.Must(session.NewSession())
		dpp.d = dynamodb.New(sess)
	}

	return dpp.d
}
