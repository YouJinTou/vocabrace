package main

import (
	"fmt"

	"github.com/YouJinTou/vocabrace/games/scrabble"
	"github.com/YouJinTou/vocabrace/pooling"
	dynamodbpooling "github.com/YouJinTou/vocabrace/pooling/providers/dynamodb"
	"github.com/google/uuid"
)

func main() {
	g := scrabble.NewGame()
	g.Print()

	for i := 0; i < 7; i++ {
		jocr := &pooling.JoinOrCreateInput{
			ConnectionID: uuid.New().String(),
			Bucket:       pooling.Novice,
			PoolLimit:    3,
		}
		provider := dynamodbpooling.NewDynamoDBProvider("dev")
		p, err := provider.JoinOrCreate(jocr)
		fmt.Println(p)
		fmt.Println(err)
		li := &pooling.LeaveInput{
			ConnectionID: jocr.ConnectionID,
			Bucket:       jocr.Bucket,
		}
		x, _ := provider.GetPool(&pooling.GetPoolInput{
			PoolID: p.ID,
			Bucket: p.Bucket,
		})
		fmt.Println(x)
		y, _ := provider.GetPeers(&pooling.GetPeersInput{
			ConnectionID: jocr.ConnectionID,
			Bucket:       p.Bucket,
		})
		fmt.Println(y)
		p2, err2 := provider.Leave(li)
		fmt.Println(p2)
		fmt.Println(err2)
	}

}
