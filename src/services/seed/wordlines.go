package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	wordlines("bulgarian", nil)
}

func wordlines(language string, startAt *string) {
	f, err := os.Open(fmt.Sprintf("resources/wordlines_%s.txt", language))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	count := 0
	totalCount := 0
	var start = false
	batch := []string{}
	s := bufio.NewScanner(f)

	for s.Scan() {
		word := strings.ToLower(s.Text())

		if !start && startAt != nil && word != *startAt {
			totalCount++
			continue
		}

		start = true

		if count < 25 {
			batch = append(batch, word)
			count++
			totalCount++
		} else {
			startTime := time.Now()
			wordlinesToDynamo(batch, totalCount, language)
			count = 0
			batch = []string{}

			wordlinesSleep(startTime)
		}
	}

	if len(batch) > 0 {
		wordlinesToDynamo(batch, totalCount, language)
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}

func wordlinesToDynamo(batch []string, totalCount int, language string) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)

	requests := []*dynamodb.WriteRequest{}
	for _, w := range batch {
		requests = append(requests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"Word": {S: aws.String(w)},
				},
			},
		})
	}

	var unprocessed map[string][]*dynamodb.WriteRequest = map[string][]*dynamodb.WriteRequest{
		fmt.Sprintf("wordlines_%s", language): requests,
	}
	for i := 0; i < 10; i++ {
		o, err := dynamo.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: unprocessed,
		})

		if len(o.UnprocessedItems) > 0 {
			unprocessed = o.UnprocessedItems
			continue
		}

		if err == nil {
			log.Printf("Count: %d. Last word written: %s", totalCount, batch[len(batch)-1])
			break
		} else {
			fmt.Printf("Capacity exceeded. Sleeping and retrying %d.", i)
			time.Sleep(time.Second * 20)
		}
	}
}

func wordlinesSleep(startTime time.Time) {
	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	if elapsed.Milliseconds() >= (time.Millisecond * 1000).Milliseconds() {
		time.Sleep(time.Second * 3)
	}
}
