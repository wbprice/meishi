package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func putRecordToTable(client *dynamodb.DynamoDB, comprehendOutput *comprehend.DetectEntitiesOutput, rawText string, prefix string) {

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	inputMap := make(map[string]*dynamodb.AttributeValue)

	for i := 0; i < len(comprehendOutput.Entities); i++ {
		entity := comprehendOutput.Entities[i]
		switch *entity.Type {
		case "PERSON":
			inputMap["name"] = &dynamodb.AttributeValue{
				S: aws.String(*entity.Text),
			}
		case "LOCATION":
			inputMap["address"] = &dynamodb.AttributeValue{
				S: aws.String(*entity.Text),
			}
		default:
			fmt.Println("Not a case comprehend understood")
		}
	}
	inputMap["raw"] = &dynamodb.AttributeValue{
		S: aws.String(rawText),
	}
	inputMap["imageLocation"] = &dynamodb.AttributeValue{
		S: aws.String(prefix),
	}

	putItemInput := dynamodb.PutItemInput{Item: inputMap, TableName: &tableName}
	_, err := client.PutItem(&putItemInput)

	if err != nil {
		log.Fatal(err)
	}
}

func handler(ctx context.Context, s3Event events.S3Event) {
	// do the thing
}

func main() {
	lambda.Start(handler)
}
