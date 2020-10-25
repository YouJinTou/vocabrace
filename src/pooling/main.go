package pool

import (
	"github.com/YouJinTou/vocabrace/memcached"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Pool handles adding and removing connections from pools.
type Pool struct {
	c *memcached.Client
}

// New creates a new pool.
func New() Pool {
	return Pool{
		c: memcached.New("localhost:11211"),
	}
}

func dynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-central-1")}))
	dynamo := dynamodb.New(sess)
	return dynamo
}
