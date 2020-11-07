package dynamodbpooling

import (
	"fmt"
	"strconv"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Leave removes a connection from the pool.
func (dpp DynamoDBProvider) Leave(i *pooling.LeaveInput) (*pooling.Pool, error) {
	for {
		c, cErr := dpp.getConnection(i.ConnectionID)

		if cErr != nil {
			panic(cErr.Error())
		}

		b, bErr := dpp.getDbBucket(i.Bucket)

		if bErr != nil {
			panic(bErr.Error())
		}

		result, err := dpp.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
			TableName: aws.String(fmt.Sprintf("%s_buckets", dpp.stage)),
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {S: aws.String(i.Bucket)},
			},
			ExpressionAttributeNames: map[string]*string{"#p": aws.String(c.PoolID)},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":cid": {SS: []*string{aws.String(c.ID)}},
				":ua":  {N: aws.String(tools.FutureTimestampStr(0))},
				":lua": {N: aws.String(strconv.Itoa(b.UpdatedAt))},
			},
			UpdateExpression:    aws.String("SET UpdatedAt = :ua DELETE #p.ConnectionIDs :cid"),
			ConditionExpression: aws.String("UpdatedAt = :lua"),
			ReturnValues:        aws.String("ALL_NEW"),
		})

		if err != nil {
			continue
		}

		dpp.detach(i.ConnectionID)

		pool := pooling.Pool{
			Bucket: i.Bucket,
		}
		dynamodbattribute.UnmarshalMap(result.Attributes[c.PoolID].M, &pool)

		return &pool, err
	}
}

func (dpp DynamoDBProvider) detach(connectionID string) {
	_, err := dpp.dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		TableName:    aws.String(fmt.Sprintf("%s_connections", dpp.stage)),
		Key:          map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(connectionID)}},
		ReturnValues: aws.String("ALL_OLD"),
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}
