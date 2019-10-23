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
)

func handler(ctx context.Context, s3Event events.S3Event) {
	session, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	inputS3BucketName := os.Getenv("BUSINESS_CARD_TEXT_S3_BUCKET_NAME")
	outputS3BucketName := os.Getenv("BUSINESS_CARD_TEXT_S3_BUCKET_NAME")

	// Configure various
	comprehendClient := comprehend.New(session)

	for i := 0; i < len(s3Event.Records); i++ {
		record := s3Event.Records[i]
		fmt.Printf("A %s event was heard.\n", record.EventName)
		fmt.Printf("The file %s was placed in the bucket %s!", record.S3.Object.Key, s3BucketName)

		inputPrefix :=
			fmt.Sprintf("s3://%s/%s", inputS3BucketName, record.S3.Object.Key)
		outputPrefix :=
			fmt.Sprintf("s3://%s/%s", inputS3BucketName, record.S3.Object.Key)

		inputDataConfig := comprehend.InputDataConfig{
			InputFormat: aws.String("ONE_DOC_PER_LINE"),
			S3Uri:       &filePrefix,
		}

			

		documentClassificationRequestInput := comprehend.StartDocumentClassificationJobInput{
			DataAccessRoleArn: 
			DocumentClassifierArn:
			InputDataConfig: &inputDataConfig,
		}

	}
}

func main() {
	lambda.Start(handler)
}
