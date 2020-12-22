package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/YouJinTou/vocabrace/facts"
	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.DynamoDBEvent) {
	for _, r := range e.Records {
		i := r.Change.NewImage
		ts, _ := strconv.Atoi(i["Timestamp"].Number())
		f := facts.Fact{
			ID:        i["ID"].String(),
			Timestamp: ts,
			Type:      i["Type"].String(),
			Data:      i["Data"].String(),
		}
		tools.SnsPublish(fmt.Sprintf("%s_events", os.Getenv("STAGE")), f)
	}
}
