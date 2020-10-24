package pool

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func dynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-central-1")}))
	dynamo := dynamodb.New(sess)
	return dynamo
}
