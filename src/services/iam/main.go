package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type user struct {
	ID       string
	Email    string
	Name     string
	Username string
}

func (u *user) hashEmail() string {
	h := sha1.New()
	h.Write([]byte(u.Email))
	s := h.Sum(nil)
	result := hex.EncodeToString(s)
	return result
}

func (u *user) getDefaultUsernameIfNotExists() string {
	if u.Username != "" {
		return u.Username
	}
	adjectives := []string{"brave", "quick", "smart", "beautiful", "agile", "sexy", "young", "strong"}
	names := []string{"badger", "ninja", "bear", "shark", "cheetah", "gorilla", "cat", "dolphin", "master"}
	rand.Seed(time.Now().UnixNano())
	username := fmt.Sprintf("%s %s",
		adjectives[rand.Intn(len(adjectives))],
		names[rand.Intn(len(names))])
	return username
}

func handler(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := providerAuth(&r)
	setResponseHeaders(&response)
	return response, err
}

func main() {
	isServerless := os.Getenv("IS_SERVERLESS") != ""

	if isServerless {
		lambda.Start(handler)
	} else {
		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowHeaders:    []string{"Content-Type"},
		}))
		r.POST("iam/provider-auth", toGinHandler(providerAuth))
		r.Run()
	}
}

func setResponseHeaders(r *events.APIGatewayProxyResponse) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers["Access-Control-Allow-Origin"] = "*"
}

func toGinHandler(f func(*events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		f(asAGWProxyRequest(c))
	}
}

func asAGWProxyRequest(c *gin.Context) *events.APIGatewayProxyRequest {
	bytes, _ := ioutil.ReadAll(c.Request.Body)
	return &events.APIGatewayProxyRequest{
		Body: string(bytes),
	}
}
