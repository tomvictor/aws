package hello

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/gin-gonic/gin"
)

func Hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello world",
	})
}

func Dynamo(ctx *gin.Context) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	// Create a DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Define the input parameters for the query
	input := &dynamodb.QueryInput{
		TableName:              aws.String("demo_table"),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: "carl"},
		},
	}

	// Execute the query and handle errors
	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		panic(err)
	}

	// Print the query results
	resp := ""
	for _, item := range result.Items {
		fmt.Println(item)
		resp = fmt.Sprintf("%v,%v", resp, item)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"result": resp,
	})
}
