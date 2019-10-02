package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, s3Event events.S3Event) {
	fmt.Printf("Hello, a file was uploaded")
}

func main() {
	lambda.Start(Handler)
}
