package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

func Handler(ctx context.Context, s3Event events.S3Event) {

	session, err := session.NewSession()
	if err != nil {
		fmt.Printf("There was an error creating a session")
	}

	s3BucketName := os.Getenv("S3_BUCKET_NAME")

	// Configure the Textract client
	textractClient := textract.New(session)

	// Iterate over file upload events
	for i := 0; i < len(s3Event.Records); i++ {
		record := s3Event.Records[i]
		fmt.Printf(record.EventName)
		fmt.Printf(record.S3.Object.Key)

		s3Object := textract.S3Object{
			Bucket: &s3BucketName,
			Name:   &record.S3.Object.Key,
		}
		document := textract.Document{
			S3Object: &s3Object,
		}
		featureTypes := aws.StringSlice([]string{"FORM"})

		analyzeDocumentInput := textract.AnalyzeDocumentInput{
			Document:     &document,
			FeatureTypes: featureTypes,
		}

		// Begin to analyze the document.
		extractOutput, err := textractClient.AnalyzeDocument(&analyzeDocumentInput)

		if err != nil {
			fmt.Printf("There was an error parsing the document")
		}
	}

}

func main() {
	lambda.Start(Handler)
}
