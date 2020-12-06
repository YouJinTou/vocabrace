package main

import (
	"encoding/json"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func providerAuth(r *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	u := &user{}
	json.Unmarshal([]byte(r.Body), u)
	if err := updateUser(u); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	uBytes, _ := json.Marshal(u)
	user := string(uBytes)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: user}, nil
}

func updateUser(u *user) error {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	u.ID = u.hashEmail()
	u.Username = u.getDefaultUsernameIfNotExists()
	ue := "SET #n = :n, Username = :u, Email = :e"
	_, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:        tools.Table("iam_users"),
		Key:              map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(u.ID)}},
		UpdateExpression: aws.String(ue),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": &dynamodb.AttributeValue{S: aws.String(u.Name)},
			":u": &dynamodb.AttributeValue{S: aws.String(u.Username)},
			":e": &dynamodb.AttributeValue{S: aws.String(u.Email)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#n": aws.String("Name"),
		},
		ReturnValues: aws.String("ALL_NEW"),
	})
	return err
}
