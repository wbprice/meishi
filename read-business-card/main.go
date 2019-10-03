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
	"github.com/aws/aws-sdk-go/service/s3"
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

	// Configure the S3 client
	s3Client := s3.New(session)

	// Iterate over file upload events
	for i := 0; i < len(s3Event.Records); i++ {
		record := s3Event.Records[i]
		fmt.Println(record.EventName)
		fmt.Println(record.S3.Object.Key)
		fmt.Println(s3BucketName)

		// Get file bytes
		blobInput := &s3.GetObjectInput{
			Bucket: aws.String(s3BucketName),
			Key:    aws.String(record.S3.Object.Key),
		}

		blob, err := s3Client.GetObject(blobInput)

		if err != nil {
			fmt.Println("Something went wrong fetching the file")
			log.Fatal(err)
		} else {
			fmt.Println("The blob")
			fmt.Println(blob)
		}

		s3Object := textract.S3Object{
			Bucket: &s3BucketName,
			Name:   &record.S3.Object.Key,
		}
		document := textract.Document{
			S3Object: &s3Object,
			// Bytes:    []byte{},
		}

		fmt.Printf("Hello %s\n", *document.S3Object.Bucket)
		fmt.Printf("I'm like %s\n", *document.S3Object.Name)

		featureTypes := aws.StringSlice([]string{"FORM"})

		analyzeDocumentInput := textract.AnalyzeDocumentInput{
			Document:     &document,
			FeatureTypes: featureTypes,
		}

		// Begin to analyze the document.
		extractOutput, err := textractClient.AnalyzeDocument(&analyzeDocumentInput)

		if err != nil {
			log.Fatal(err)
		} else {
			metadata := extractOutput.DocumentMetadata
			fmt.Printf("%d pages were found", metadata.Pages)
		}
	}
}

func main() {
	lambda.Start(handler)
}
