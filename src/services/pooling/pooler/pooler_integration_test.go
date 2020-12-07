package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/YouJinTou/vocabrace/services/pooling/ws"
	"github.com/YouJinTou/vocabrace/tools"
)

var userID = "test_user_id"
var mappings = func() []*mapping {
	return []*mapping{
		&mapping{ConnectionID: r(), UserID: r()},
		&mapping{ConnectionID: r(), UserID: r()},
		&mapping{ConnectionID: r(), UserID: r()},
		&mapping{ConnectionID: r(), UserID: r()},
		&mapping{ConnectionID: r(), UserID: r()},
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func Test_OnConnect_JoinsWaitlist(t *testing.T) {
	if _, err := p().joinWaitlist(r(), params(r(), r(), r(), r())); err != nil {
		t.Errorf(err.Error())
	}
}

func Test_OnConnect_CreatesBucket(t *testing.T) {
	params := params(r(), r(), r(), r())
	p().joinWaitlist(r(), params)
	bucket, err := bucket(params)

	if bucket.Item == nil {
		t.Errorf("key not found")
	}
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_OnConnect_SetsConnectionIDs(t *testing.T) {
	params := params(r(), r(), r(), r())
	mappings := mappings()
	for _, m := range mappings {
		p().joinWaitlist(m.ConnectionID, params)
	}
	bucket, _ := bucket(params)

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
	params := params(r(), r(), r(), r())
	mappings := mappings()
	for _, m := range mappings {
		params["userID"] = m.UserID
		p().joinWaitlist(m.ConnectionID, params)
	}
	bucket, _ := bucket(params)

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
	params := params(r(), r(), r(), r())
	mappings := mappings()
	for _, m := range mappings {
		params["userID"] = m.UserID
		p().joinWaitlist(m.ConnectionID, params)
	}
	params["game"] = ""
	err := p().leaveWaitlist(mappings[1].ConnectionID, params)

	if err == nil {
		bucket, _ := bucket(params)
		connectionIDs := bucket.Item["ConnectionIDs"].SS
		if tools.ContainsStringPtr(connectionIDs, mappings[1].ConnectionID) {
			t.Errorf("found connection, expected to be removed")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_OnDisconnect_RemovesMapping(t *testing.T) {
	params := params(r(), r(), r(), r())
	mappings := mappings()
	for _, m := range mappings {
		params["userID"] = m.UserID
		p().joinWaitlist(m.ConnectionID, params)
	}
	params["game"] = ""
	err := p().leaveWaitlist(mappings[4].ConnectionID, params)

	if err == nil {
		bucket, _ := bucket(params)
		if tools.ContainsStringPtr(bucket.Item["Mappings"].SS, mappings[4].String()) {
			t.Errorf("expected %s to have been removed", mappings[4].String())
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_OnWaitlistNotFull_DoesNotFlush(t *testing.T) {
	var oPtr *dynamodb.UpdateItemOutput
	params := params(r(), r(), r(), r())
	for _, m := range mappings()[0:3] {
		params["userID"] = m.UserID
		o, _ := p().joinWaitlist(m.ConnectionID, params)
		oPtr = o
	}

	p().onWaitlistFull(oPtr, req(params))

	bucket, _ := tools.GetItem(tools.Table("waitlist"), "ID", ID(params), nil, nil, nil)

	if _, ok := bucket.Item["ConnectionIDs"]; !ok {
		t.Errorf("expected connections to be there")
	}
	if _, ok := bucket.Item["Mappings"]; !ok {
		t.Errorf("expected mappings to be there")
	}
}

func Test_OnWaitlistFull_FlushesWaitlist(t *testing.T) {
	var oPtr *dynamodb.UpdateItemOutput
	mappings := mappings()
	params := params(r(), r(), strconv.Itoa(len(mappings)), r())
	p := p()
	for _, m := range mappings {
		params["userID"] = m.UserID
		o, _ := p.joinWaitlist(m.ConnectionID, params)
		oPtr = o
	}

	p.onWaitlistFull(oPtr, req(params))

	bucket, _ := tools.GetItem(tools.Table("waitlist"), "ID", ID(params), nil, nil, nil)

	if _, ok := bucket.Item["ConnectionIDs"]; ok {
		t.Errorf("expected connections to have been cleared")
	}
	if _, ok := bucket.Item["Mappings"]; ok {
		t.Errorf("expected mappings to have been cleared")
	}
}

func Test_OnWaitlistFull_PoolsPlayers(t *testing.T) {
	var oPtr *dynamodb.UpdateItemOutput
	mappings := mappings()
	params := params(r(), r(), strconv.Itoa(len(mappings)), r())
	for _, m := range mappings {
		params["userID"] = m.UserID
		o, _ := p().joinWaitlist(m.ConnectionID, params)
		oPtr = o
	}

	p().onWaitlistFull(oPtr, req(params))
}

func Test_OnWaitlistFull_CreatesPool(t *testing.T) {
	var oPtr *dynamodb.UpdateItemOutput
	var pid string
	mappings := mappings()
	params := params(r(), r(), strconv.Itoa(len(mappings)), r())
	p := pos(func(i *ws.OnStartInput) { pid = i.PoolID })
	for _, m := range mappings {
		params["userID"] = m.UserID
		o, _ := p.joinWaitlist(m.ConnectionID, params)
		oPtr = o
	}

	p.onWaitlistFull(oPtr, req(params))

	if result, err := tools.GetItem(tools.Table("pools"), "ID", pid, nil, nil, nil); err == nil {
		if result.Item == nil {
			t.Errorf("pool not created")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_PlayersPooled_SetsConnectionsData(t *testing.T) {
	var pid string
	var oPtr *dynamodb.UpdateItemOutput
	mappings := mappings()
	params := params(r(), r(), strconv.Itoa(len(mappings)), r())
	p := pos(func(i *ws.OnStartInput) { pid = i.PoolID })
	for _, m := range mappings {
		params["userID"] = m.UserID
		o, _ := p.joinWaitlist(m.ConnectionID, params)
		oPtr = o
	}

	p.onWaitlistFull(oPtr, req(params))

	connection, err := tools.GetItem(
		tools.Table("connections"), "ID", mappings[1].ConnectionID, nil, nil, nil)
	if err == nil {
		if val, ok := connection.Item["PoolID"]; ok {
			if val.S == nil || *val.S != pid {
				t.Errorf("incorrect pool ID set")
			}
		} else {
			t.Errorf("pool ID not set for connection")
		}
		if val, ok := connection.Item["Waitlist"]; ok {
			if val.S == nil || *val.S != p.getBucket(params) {
				t.Errorf("incorrect waitlist set")
			}
		} else {
			t.Errorf("waitlist not set for connection")
		}
		if val, ok := connection.Item["UserID"]; ok {
			if val.S == nil || *val.S != mappings[1].UserID {
				t.Errorf("incorrect user ID set")
			}
		} else {
			t.Errorf("user ID not set for connection")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func req(params map[string]string) *events.APIGatewayWebsocketProxyRequest {
	return &events.APIGatewayWebsocketProxyRequest{
		QueryStringParameters: params,
		RequestContext: events.APIGatewayWebsocketProxyRequestContext{
			DomainName:   "domain",
			ConnectionID: "cid",
		},
	}
}

func params(g, l, p, u string) map[string]string {
	return map[string]string{
		"game":     g,
		"language": l,
		"players":  p,
		"userID":   u,
	}
}

func bucket(params map[string]string) (*dynamodb.GetItemOutput, error) {
	return tools.GetItem(tools.Table("waitlist"), "ID", ID(params), nil, nil, nil)
}

func ID(params map[string]string) string {
	return p().getBucket(params)
}

func setup() {
	os.Setenv("STAGE", "dev")
	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("AWS_PROFILE", "vocabrace")
}

func p() *pooler {
	return &pooler{
		OnStart: func(i *ws.OnStartInput) {},
	}
}

func pos(f func(*ws.OnStartInput)) pooler {
	return pooler{OnStart: f}
}

func r() string {
	return uuid.New().String()[0:5]
}
