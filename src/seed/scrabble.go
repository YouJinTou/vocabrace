package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	scrabble("bulgarian")
}

func scrabble(language string) {
	f, err := os.Open(fmt.Sprintf("%s.txt", language))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	count := 0
	totalCount := 0
	batch := []string{}
	for s.Scan() {
		if count < 100 {
			batch = append(batch, s.Text())
			count++
			totalCount++
		} else {
			scrabbleToDynamo(batch, totalCount, language)
			count = 0
			batch = []string{}
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}

func scrabbleToDynamo(batch []string, totalCount int, language string) {
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
	_, err := dynamo.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			fmt.Sprintf("scrabble_%s", language): requests,
		},
	})

	if err == nil {
		log.Printf("Count: %d. Last word written: %s", totalCount, batch[len(batch)-1])
	} else {
		log.Printf("Failed to process batch. Reason: %s", err.Error())
	}
}
