package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/iam/provider-auth", providerAuth)
	r.Run()
}

func providerAuth(c *gin.Context) {
	c.JSON(200, gin.H{
		"result": 1,
	})
}
