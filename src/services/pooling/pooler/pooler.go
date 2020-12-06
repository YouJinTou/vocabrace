package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/YouJinTou/vocabrace/services/pooling/ws"
	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type mapping struct {
	ConnectionID string
	UserID       string
}

func (m *mapping) String() string {
	return fmt.Sprintf("%s|%s", m.ConnectionID, m.UserID)
}

func toMapping(s string) mapping {
	tokens := strings.Split(s, "|")
	return mapping{ConnectionID: tokens[0], UserID: tokens[1]}
}

type pooler struct {
	OnStart func(*ws.ReceiverData)
}

func (p *pooler) joinWaitlist(connectionID string, params map[string]string) (
	*dynamodb.UpdateItemOutput, error) {
	mapping := p.getMapping(connectionID, params)
	i := dynamodb.UpdateItemInput{
		TableName: table("waitlist"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(p.getBucket(params))}},
		UpdateExpression: aws.String("ADD ConnectionIDs :c, Mappings :m"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": &dynamodb.AttributeValue{SS: []*string{aws.String(connectionID)}},
			":m": &dynamodb.AttributeValue{SS: []*string{aws.String(mapping)}},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	o, err := dynamo().UpdateItem(&i)
	return o, err
}

func (p *pooler) leaveWaitlist(connectionID string, params map[string]string) error {
	mapping := p.getMapping(connectionID, params)
	i := dynamodb.UpdateItemInput{
		TableName: table("waitlist"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(p.getBucket(params))}},
		UpdateExpression: aws.String("DELETE ConnectionIDs :c, Mappings :m"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": &dynamodb.AttributeValue{SS: []*string{aws.String(connectionID)}},
			":m": &dynamodb.AttributeValue{SS: []*string{aws.String(mapping)}},
		},
	}
	_, err := dynamo().UpdateItem(&i)
	return err
}

func (p *pooler) onWaitlistFull(o *dynamodb.UpdateItemOutput, r *events.APIGatewayWebsocketProxyRequest) {
	players, _ := strconv.Atoi(getParam("players", r.QueryStringParameters))
	connectionIDs := tools.FromStringPtrs(o.Attributes["ConnectionIDs"].SS)
	poolFull := len(connectionIDs) == players
	if !poolFull {
		return
	}

	users := []*ws.User{}
	for _, s := range o.Attributes["Mappings"].SS {
		m := toMapping(*s)
		users = append(users, &ws.User{
			ConnectionID: m.ConnectionID,
			Username:     "seom-test",
			UserID:       m.UserID,
		})
	}

	pid := p.createPool(connectionIDs)
	p.setPoolForConnections(connectionIDs, pid)

	p.OnStart(&ws.ReceiverData{
		Users:         users,
		PoolID:        pid,
		ConnectionIDs: connectionIDs,
		Domain:        r.RequestContext.DomainName,
		Stage:         os.Getenv("STAGE"),
		Game:          getParam("game", r.QueryStringParameters),
		Language:      getParam("language", r.QueryStringParameters),
	})

	p.flushWaitlist(r.QueryStringParameters)
}

func (p *pooler) flushWaitlist(params map[string]string) error {
	i := dynamodb.UpdateItemInput{
		TableName: table("waitlist"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(p.getBucket(params))}},
		UpdateExpression: aws.String("REMOVE ConnectionIDs, Mappings"),
	}
	_, err := dynamo().UpdateItem(&i)
	return err
}

func (p *pooler) createPool(connectionIDs []string) string {
	ID := uuid.New().String()
	tools.PutItem(*table("pools"), struct {
		ID            string
		ConnectionIDs []string
	}{ID, connectionIDs})
	return ID
}

func (p *pooler) setPoolForConnections(connectionIDs []string, poolID string) {
	for _, cid := range connectionIDs {
		tools.PutItem(*table("connections"), struct {
			ID        string
			PoolID    string
			LiveUntil int
		}{cid, poolID, tools.FutureTimestamp(7200)})
	}
}

func (p *pooler) getSkill(params map[string]string) string {
	return "novice"
}

func (p *pooler) getBucket(params map[string]string) string {
	bucket := fmt.Sprintf("%s_%s_%s_%s_%s",
		os.Getenv("STAGE"),
		getParam("game", params),
		getParam("language", params),
		p.getSkill(params),
		getParam("players", params))
	return bucket
}

func (p *pooler) getMapping(connectionID string, params map[string]string) string {
	mapping := mapping{
		ConnectionID: connectionID,
		UserID:       uuid.New().String(),
	}
	isAnonymousVal := getNilParam("isAnonymous", params)
	isAnonymous := isAnonymousVal == nil || *isAnonymousVal == "true"
	userID := getNilParam("userID", params)
	if isAnonymous && userID == nil {
		return mapping.String()
	}
	mapping.UserID = *userID
	return mapping.String()
}

func getNilParam(key string, params map[string]string) *string {
	if val, ok := params[key]; ok {
		return &val
	}
	return nil
}

func getParam(key string, params map[string]string) string {
	if val, ok := params[key]; ok {
		return val
	}
	panic(fmt.Sprintf("expected param %s", key))
}

func dynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}

func table(name string) *string {
	return aws.String(fmt.Sprintf("%s_%s", os.Getenv("STAGE"), name))
}
