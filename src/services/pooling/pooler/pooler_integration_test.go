package main

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/YouJinTou/vocabrace/services/pooling/ws"
	"github.com/YouJinTou/vocabrace/tools"
)

var userID = "test_user_id"
var mappings = []*mapping{
	&mapping{ConnectionID: "123", UserID: "qqq"},
	&mapping{ConnectionID: "456", UserID: "www"},
	&mapping{ConnectionID: "789", UserID: "eee"},
	&mapping{ConnectionID: "999", UserID: "rrr"},
	&mapping{ConnectionID: "000", UserID: "ttt"},
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func Test_OnConnect_JoinsWaitlist(t *testing.T) {
	if _, err := p().joinWaitlist("123", params(userID)); err != nil {
		t.Errorf(err.Error())
	}
}

func Test_OnConnect_CreatesBucket(t *testing.T) {
	p().joinWaitlist("123", params(userID))
	bucket, err := bucket()

	if bucket.Item == nil {
		t.Errorf("key not found")
	}
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_OnConnect_SetsConnectionIDs(t *testing.T) {
	for _, m := range mappings {
		p().joinWaitlist(m.ConnectionID, params(m.UserID))
	}
	bucket, _ := bucket()

	if value, ok := bucket.Item["ConnectionIDs"]; ok {
		for _, m := range mappings {
			if !tools.ContainsStringPtr(value.SS, m.ConnectionID) {
				t.Errorf("%s not found", m.ConnectionID)
			}
		}
	} else {
		t.Errorf("'ConnectionIDs' not set")
	}
}

func Test_OnConnect_SetsConnectionIDUserIDMapping(t *testing.T) {
	for _, m := range mappings {
		p().joinWaitlist(m.ConnectionID, params(m.UserID))
	}
	bucket, _ := bucket()

	if value, ok := bucket.Item["Mappings"]; ok {
		for _, m := range mappings {
			if !tools.ContainsStringPtr(value.SS, m.String()) {
				t.Errorf("%s not found", m.String())
			}
		}
	} else {
		t.Errorf("'Mappings' not set")
	}
}

func Test_OnDisconnect_RemovesConnectionsFromWaitlist(t *testing.T) {
	for _, m := range mappings {
		p().joinWaitlist(m.ConnectionID, params(m.UserID))
	}
	err := p().leaveWaitlist(mappings[1].ConnectionID, params(mappings[1].UserID))

	if err == nil {
		bucket, _ := bucket()
		connectionIDs := bucket.Item["ConnectionIDs"].SS
		if tools.ContainsStringPtr(connectionIDs, mappings[1].ConnectionID) {
			t.Errorf("found connection, expected to be removed")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_OnDisconnect_RemovesMapping(t *testing.T) {
	for _, m := range mappings {
		p().joinWaitlist(m.ConnectionID, params(m.UserID))
	}
	err := p().leaveWaitlist(mappings[1].ConnectionID, params(mappings[1].UserID))

	if err == nil {
		bucket, _ := bucket()
		if tools.ContainsStringPtr(bucket.Item["Mappings"].SS, mappings[1].String()) {
			t.Errorf("expected %s to have been removed", mappings[1].String())
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_OnWaitlistNotFull_DoesNotFlush(t *testing.T) {
	var oPtr *dynamodb.UpdateItemOutput
	for _, m := range mappings[0:3] {
		o, _ := p().joinWaitlist(m.ConnectionID, params(m.UserID))
		oPtr = o
	}

	p().onWaitlistFull(oPtr, req(userID))

	bucket, _ := tools.GetItem("dev_waitlist", "ID", ID(), nil, nil)

	if _, ok := bucket.Item["ConnectionIDs"]; !ok {
		t.Errorf("expected connections to be there")
	}
	if _, ok := bucket.Item["Mappings"]; !ok {
		t.Errorf("expected mappings to be there")
	}
}

func Test_OnWaitlistFull_FlushesWaitlist(t *testing.T) {
	var oPtr *dynamodb.UpdateItemOutput
	for _, m := range mappings {
		o, _ := p().joinWaitlist(m.ConnectionID, params(m.UserID))
		oPtr = o
	}

	p().onWaitlistFull(oPtr, req(userID))

	bucket, _ := tools.GetItem("dev_waitlist", "ID", ID(), nil, nil)

	if _, ok := bucket.Item["ConnectionIDs"]; ok {
		t.Errorf("expected connections to have been cleared")
	}
	if _, ok := bucket.Item["Mappings"]; ok {
		t.Errorf("expected mappings to have been cleared")
	}
}

func Test_OnWaitlistFull_PoolsPlayers(t *testing.T) {
	var oPtr *dynamodb.UpdateItemOutput
	for _, m := range mappings {
		o, _ := p().joinWaitlist(m.ConnectionID, params(m.UserID))
		oPtr = o
	}

	p().onWaitlistFull(oPtr, req(userID))
}

func Test_PlayersPooled_SetsConnectionPoolMappings(t *testing.T) {
	var oPtr *dynamodb.UpdateItemOutput
	for _, m := range mappings {
		o, _ := p().joinWaitlist(m.ConnectionID, params(m.UserID))
		oPtr = o
	}

	p().onWaitlistFull(oPtr, req(userID))

	connection, err := tools.GetItem("dev_connections", "ID", mappings[1].ConnectionID, nil, nil)
	if err == nil {
		if val, ok := connection.Item["PoolID"]; ok {
			if val.S == nil || *val.S != "test_pool_id" {
				t.Errorf("incorrect pool ID set")
			}
		} else {
			t.Errorf("pool ID not set for connection")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func req(userID string) *events.APIGatewayWebsocketProxyRequest {
	return &events.APIGatewayWebsocketProxyRequest{
		QueryStringParameters: params(userID),
		RequestContext: events.APIGatewayWebsocketProxyRequestContext{
			DomainName:   "domain",
			ConnectionID: "cid",
		},
	}
}

func params(userID string) map[string]string {
	return map[string]string{
		"game":     "scrabble",
		"language": "bulgarian",
		"players":  "5",
		"userID":   userID,
	}
}

func bucket() (*dynamodb.GetItemOutput, error) {
	return tools.GetItem("dev_waitlist", "ID", ID(), nil, nil)
}

func ID() string {
	return p().getBucket(params(userID))
}

func setup() {
	os.Setenv("STAGE", "dev")
	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("AWS_PROFILE", "vocabrace")
}

func shutdown() {
	tools.DeleteItem("dev_waitlist", "ID", ID(), nil, nil)
}

func p() *pooler {
	return &pooler{
		OnStart: func(*ws.ReceiverData) ws.PoolID {
			return "test_pool_id"
		},
	}
}
