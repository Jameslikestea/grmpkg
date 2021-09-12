package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"grmpkg.com/grmpkg/internal/model"
)

func main() {
	sess := session.New(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	ddb := dynamodb.New(sess)

	item, err := dynamodbattribute.MarshalMap(model.Package{
		Name: "e9e675eb5530.ngrok.io/some/other",
		Info: model.PackageInfo{
			Package:  "some/other",
			Hostname: "e9e675eb5530.ngrok.io",
		},
		Versions: []model.PackageVersion{
			{
				Name:    "v1.0.0",
				Short:   "v1.0.0",
				Version: "v1.0.0",
				Time:    time.Now().Add(time.Hour * -4).Format(time.RFC3339),
			},
			{
				Name:    "v1.1.0",
				Short:   "v1.1.0",
				Version: "v1.1.0",
				Time:    time.Now().Add(time.Hour * -2).Format(time.RFC3339),
			},
			{
				Name:    "v1.2.0",
				Short:   "v1.2.0",
				Version: "v1.2.0",
				Time:    time.Now().Add(time.Hour * -1).Format(time.RFC3339),
			},
			{
				Name:    "v1.2.1",
				Short:   "v1.2.1",
				Version: "v1.2.1",
				Time:    time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = ddb.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("grmpkg"),
		Item:      item,
	})
	if err != nil {
		panic(err)
	}
}
