package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"log"
)

func main() {
	fmt.Println("List aws cognito pools")

	ctx := context.Background()

	defaultConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal(err.Error())
	}

	client := cognitoidentityprovider.NewFromConfig(defaultConfig)
	userPoolInput := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: 10,
	}
	pools, err := client.ListUserPools(ctx, userPoolInput)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, pool := range pools.UserPools {
		fmt.Println(*pool.Name)
	}

}
