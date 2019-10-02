package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

func Handler(ctx context.Context, s3Event events.S3Event) {

	session, err := session.NewSession()
	if err != nil {
		fmt.Printf("There was an error creating a session")
	}

	// Configure the Textract client
	textract_client := textract.New(session)

	// Iterate over file upload events
	for i := 0; i < len(s3Event.Records); i++ {
		record := s3Event.Records[i]
		fmt.Printf(record.EventName)
		fmt.Printf(record.S3.Object.Key)

		// Begin to analyze the document.
		extract_output, err := textract_client.AnalyzeDocument(
			AnalyzeDocumentInput{
				Document{
					S3Object: record.S3.Object,
				},
				["FORMS"]
			},
		)

		if err != nil {
			fmt.Printf("There was an error parsing the document")
		}
	}

}

func main() {
	lambda.Start(Handler)
}
