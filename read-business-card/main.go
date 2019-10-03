package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	// Does the thing

	session, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	s3BucketName := os.Getenv("S3_BUCKET_NAME")

	// Configure the Textract client
	textractClient := textract.New(session)

	// Iterate over file upload events
	for i := 0; i < len(s3Event.Records); i++ {
		record := s3Event.Records[i]
		fmt.Printf("A %s event was heard.\n", record.EventName)
		fmt.Printf("The file %s was placed in the bucket %s!", record.S3.Object.Key, s3BucketName)

		if err != nil {
			fmt.Println("Something went wrong fetching the file")
			log.Fatal(err)
		}

		// Create an S3Object to use with texttract.Document
		s3Object := textract.S3Object{
			Bucket: aws.String(s3BucketName),
			Name:   aws.String(record.S3.Object.Key),
		}

		// Create the input for the call to textractClient.DetectDocumentText
		document := textract.Document{
			S3Object: &s3Object,
		}
		detectDocumentTextInput := textract.DetectDocumentTextInput{
			Document: &document,
		}

		// Begin to analyze the document.
		extractOutput, err := textractClient.DetectDocumentText(&detectDocumentTextInput)

		if err != nil {
			log.Fatal(err)
		} else {
			jsonBytes, err := json.MarshalIndent(extractOutput.GoString(), "", "    ")
			if err != nil {
				fmt.Println("There was an issue formatting the response.")
				log.Fatal(err)
			}
			fmt.Printf("%s", jsonBytes)
			metadata := extractOutput.DocumentMetadata
			fmt.Printf("%d pages were found", metadata.Pages)
		}
	}
}

func main() {
	lambda.Start(handler)
}
