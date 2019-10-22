package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/textract"
)

func analyzeBusinessCardText(client *comprehend.Comprehend, text *string) *comprehend.DetectEntitiesOutput {
	// Packages text in a call to AWS Comprehend

	languageCode := "en"
	detectEntitiesInput := comprehend.DetectEntitiesInput{
		LanguageCode: &languageCode,
		Text:         text,
	}

	comprehendOutput, err := client.DetectEntities(&detectEntitiesInput)

	if err != nil {
		log.Fatal(err)
	}

	return comprehendOutput
}

func sortBusinessCardText(comprehendOutput *comprehend.DetectEntitiesOutput) {
	for i := 0; i < len(comprehendOutput.Entities); i++ {
		entity := comprehendOutput.Entities[i]
		fmt.Printf("String: %s\n", *entity.Text)
		fmt.Printf("Type: %s\n", *entity.Type)
	}
}

func putRecordToTable(client *dynamodb.DynamoDB, comprehendOutput *comprehend.DetectEntitiesOutput, s3Object *textract.S3Object) {

}

func getTextFromBusinessCard(client *textract.Textract, s3Object textract.S3Object) *textract.DetectDocumentTextOutput {
	// Create the input for the call to textractClient.DetectDocumentText
	document := textract.Document{
		S3Object: &s3Object,
	}
	detectDocumentTextInput := textract.DetectDocumentTextInput{
		Document: &document,
	}

	// Begin to analyze the document.
	extractOutput, err := client.DetectDocumentText(&detectDocumentTextInput)

	if err != nil {
		log.Fatal(err)
	}

	return extractOutput
}

func flattenTextFromTextractOutputBlocks(extractOutput *textract.DetectDocumentTextOutput) *string {
	output := ""
	for i := 0; i < len(extractOutput.Blocks); i++ {
		block := extractOutput.Blocks[i]

		if *block.BlockType == "LINE" {
			output += *block.Text + "\n"
		}
	}
	return &output
}

func handler(ctx context.Context, s3Event events.S3Event) {
	// Does the thing

	session, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	s3BucketName := os.Getenv("S3_BUCKET_NAME")

	// Configure various service clients
	textractClient := textract.New(session)
	comprehendClient := comprehend.New(session)
	dynamoDbClient := dynamodb.New(session)

	// Iterate over file upload events
	for i := 0; i < len(s3Event.Records); i++ {
		record := s3Event.Records[i]
		fmt.Printf("A %s event was heard.\n", record.EventName)
		fmt.Printf("The file %s was placed in the bucket %s!", record.S3.Object.Key, s3BucketName)

		// Create an S3Object to use with texttract.Document
		s3Object := textract.S3Object{
			Bucket: aws.String(s3BucketName),
			Name:   aws.String(record.S3.Object.Key),
		}
		// Get analysis of image from Textract.
		documentOutput := getTextFromBusinessCard(textractClient, s3Object)
		documentText := flattenTextFromTextractOutputBlocks(documentOutput)

		// Get interesting lines of text from documentOutput
		comprehendOutput := analyzeBusinessCardText(comprehendClient, documentText)
		// Look at each line
		sortBusinessCardText(comprehendOutput)
		// Save record to table
		putRecordToTable(dynamoDbClient, comprehendOutput)
	}
}

func main() {
	lambda.Start(handler)
}
