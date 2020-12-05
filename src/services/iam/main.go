package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// User is the application user.
type User struct {
	ID       string
	Email    string
	Name     string
	Username string
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
