package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	S3Client     *s3.Client
	S3BucketName string
)

// InitS3 initializes the S3 client and loads credentials/configs
func InitS3() {
	S3BucketName = os.Getenv("S3_BUCKET_NAME")
	if S3BucketName == "" {
		log.Fatal("S3_BUCKET_NAME environment variable is not set")
	}

	ctx := context.Background()
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Fatal("AWS_REGION environment variable is not set")
	}

	// ✅ Pass context first, then options
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	S3Client = s3.NewFromConfig(cfg)
	log.Println("Connected to AWS S3 successfully")
}
