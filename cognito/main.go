package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
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
	firstPool := DescribeCognitoUserPools(err, client, ctx)
	DescribeUsers(firstPool, err, client, ctx)

	// get access token
	//client.InitiateAuth()
}

func DescribeUsers(firstPool types.UserPoolDescriptionType, err error, client *cognitoidentityprovider.Client, ctx context.Context) {
	userInput := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: firstPool.Id,
	}
	users, err := client.ListUsers(ctx, userInput)

	fmt.Println(len(users.Users))

	for _, user := range users.Users {
		fmt.Println(*user.Username)
		for _, att := range user.Attributes {
			fmt.Println(*att.Name, *att.Value)
		}
	}
}

func DescribeCognitoUserPools(err error, client *cognitoidentityprovider.Client, ctx context.Context) types.UserPoolDescriptionType {
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

	return pools.UserPools[0]
}
