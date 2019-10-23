package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/textract"
)

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

func putTextToS3(client *s3.S3, key *string, text *string) {
	buffer := []byte(*text)
	businessCardTextS3BucketName := os.Getenv("BUSINESS_CARD_TEXT_S3_BUCKET_NAME")

	putObjectInput := s3.PutObjectInput{
		Bucket:             &businessCardTextS3BucketName,
		Key:                key,
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(int64(len(buffer))),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	}

	_, err := client.PutObject(&putObjectInput)

	if err != nil {
		log.Fatal(err)
	}
}

func handler(ctx context.Context, s3Event events.S3Event) {
	// Does the thing

	session, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	businessCardImageS3BucketName := os.Getenv("BUSINESS_CARD_IMAGE_S3_BUCKET_NAME")

	// Configure the textract client
	textractClient := textract.New(session)
	s3Client := s3.New(session)

	// Iterate over file upload events
	for i := 0; i < len(s3Event.Records); i++ {
		record := s3Event.Records[i]
		fmt.Printf("A %s event was heard.\n", record.EventName)
		fmt.Printf("The file %s was placed in the bucket %s!", record.S3.Object.Key, businessCardImageS3BucketName)

		// Create an S3Object to use with texttract.Document
		s3Object := textract.S3Object{
			Bucket: aws.String(businessCardImageS3BucketName),
			Name:   aws.String(record.S3.Object.Key),
		}
		// Get analysis of image from Textract.
		documentOutput := getTextFromBusinessCard(textractClient, s3Object)
		documentText := flattenTextFromTextractOutputBlocks(documentOutput)

		// Save document text to a CSV in the destination bucket
		putTextToS3(s3Client, &record.S3.Object.Key, documentText)
	}
}

func main() {
	lambda.Start(handler)
}
