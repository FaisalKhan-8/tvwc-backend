package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Client represents the AWS S3 client.
type S3Client struct {
	Client *s3.Client
	Bucket string
}

// NewS3Client configures and returns a new S3 client.
func NewS3Client() (*S3Client, error) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_BUCKET")

	if accessKeyID == "" || secretAccessKey == "" || region == "" || bucket == "" {
		return nil, errors.New("missing AWS credentials or bucket information in environment variables")
	}

	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Client{
		Client: client,
		Bucket: bucket,
	}, nil
}

// UploadFile uploads a file to S3.
func (s *S3Client) UploadFile(ctx context.Context, key string, body io.Reader) (*s3.PutObjectOutput, error) {
	_, err := s.Client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(s.Bucket)})
    if err != nil {
        var notFoundError *types.NotFound
        if errors.As(err, &notFoundError) {
            log.Printf("Bucket %s does not exist or is not accessible. Creating...", s.Bucket)
            _, err := s.Client.CreateBucket(ctx, &s3.CreateBucketInput{
                Bucket: aws.String(s.Bucket),
                CreateBucketConfiguration: &types.CreateBucketConfiguration{
                    LocationConstraint: types.BucketLocationConstraint(os.Getenv("AWS_REGION")),
                },
            })
            if err != nil {
                return nil, fmt.Errorf("error creating bucket: %w", err)
            }
            log.Printf("Bucket %s created successfully.", s.Bucket)
        } else {
            return nil, fmt.Errorf("error checking bucket: %w", err)
        }
    }

	input := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	output, err := s.Client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return output, nil
}

// DeleteFile deletes a file from S3.
func (s *S3Client) DeleteFile(ctx context.Context, key string) (*s3.DeleteObjectOutput, error) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	output, err := s.Client.DeleteObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	return output, nil
}

// ListFiles lists files in an S3 bucket.
func (s *S3Client) ListFiles(ctx context.Context, prefix string) ([]types.Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(prefix),
	}

	output, err := s.Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return output.Contents, nil
}

// Example of how to use the UploadFile function.
func main() {
	client, err := NewS3Client()
	if err != nil {
		log.Fatalf("Error creating S3 client: %v", err)
	}

	file, err := os.Open("test.txt") // Replace with your file path
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output, err := client.UploadFile(ctx, "test.txt", file) // Replace with desired key
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	log.Printf("File uploaded successfully: %v", output)
}