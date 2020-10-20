package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/api/v1/healthcheck", healthcheck)
	router.Run()
}

func healthcheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": 200,
	})
}
