package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Initialize Gin router
	r := gin.Default()

	// Set up your routes
	r.Any("/", func(c *gin.Context) {
		name := c.DefaultQuery("name", "world")
		message := fmt.Sprintf("Hello, %s!!!", name)
		c.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	})

	// Create the Lambda handler
	ginLambda = ginadapter.New(r)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Proxy the API Gateway request to Gin
	return ginLambda.Proxy(req)
}

func main() {
	// Start the Lambda handler
	lambda.Start(Handler)
}
