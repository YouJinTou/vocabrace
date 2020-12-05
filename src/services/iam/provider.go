package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

func providerAuth(r *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	u := &User{}
	json.Unmarshal([]byte(r.Body), u)
	if err := updateUser(u); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	uBytes, _ := json.Marshal(u)
	user := string(uBytes)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: user}, nil
}

func updateUser(u *User) error {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	ue := "SET #n = :n, Username = if_not_exists(Username, :u), ID = if_not_exists(ID, :i)"
	o, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:        aws.String(fmt.Sprintf("%s_iam_users", os.Getenv("STAGE"))),
		Key:              map[string]*dynamodb.AttributeValue{"Email": {S: aws.String(u.Email)}},
		UpdateExpression: aws.String(ue),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": &dynamodb.AttributeValue{S: aws.String(u.Name)},
			":u": &dynamodb.AttributeValue{S: aws.String(generateUsername())},
			":i": &dynamodb.AttributeValue{S: aws.String(uuid.New().String())},
		},
		ExpressionAttributeNames: map[string]*string{
			"#n": aws.String("Name"),
		},
		ReturnValues: aws.String("ALL_NEW"),
	})
	if err == nil {
		u.ID = *o.Attributes["ID"].S
		u.Username = *o.Attributes["Username"].S
	}
	return err
}

func generateUsername() string {
	id := uuid.New().String()
	candidate := strings.Replace(id, "-", "", -1)
	username := candidate[0:6]
	return username
}
