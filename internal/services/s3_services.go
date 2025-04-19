package services

import (
	"Backend/configs"
	"bytes"
	"context"
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	s3Client *s3.Client
	bucket   string
}

func NewAWSService() (*S3Service, error) {
	s3Config := configs.LoadConfig()
	var region = s3Config.AWSRegion
	var bucket = s3Config.S3Bucket
	var accessKey = s3Config.AWSAccessKeyId
	var secretKey = s3Config.AWSSecretAccessKey

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))
	if err != nil {
		return nil, err
	}

	// Create an Amazon S3 service client access key and so on
	s3Client := s3.NewFromConfig(cfg)

	return &S3Service{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func NewR2Service() (*S3Service, error) {
	s3Config := configs.LoadConfig()
	// Ensure bucket name is not empty
	var bucket = s3Config.S3Bucket
	if bucket == "" {
		// Use the correct bucket name from Cloudflare R2
		bucket = "pufa-2025" // Correct bucket name from Cloudflare dashboard
	}
	
	var accessKey = s3Config.CloudflareR2AccessId
	var secretKey = s3Config.CloudflareR2AccessKey
	var url = "https://" + s3Config.CloudflareAccountId + ".r2.cloudflarestorage.com/"

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: url,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("apac"),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &S3Service{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func (s *S3Service) UploadFileToAWS(ctx context.Context, directory, key string, file []byte) error {
	key = directory + "/" + key + ".jpg"

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/jpeg"),
	}

	_, err := s.s3Client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) UploadFileToR2(ctx context.Context, directory, key string, file []byte) error {
	// Ensure we have a valid bucket name
	if s.bucket == "" {
		return fmt.Errorf("bucket name is empty")
	}
	
	// First, check if the bucket exists
	fmt.Printf("Checking if bucket exists: %s\n", s.bucket)
	_, err := s.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	
	if err != nil {
		// If the bucket doesn't exist, log detailed error info
		fmt.Printf("Error checking bucket existence: %v\n", err)
		fmt.Printf("Cloudflare account ID: %s\n", os.Getenv("CLOUDFLARE_ACCOUNT_ID"))
		fmt.Printf("R2 access key ID length: %d\n", len(os.Getenv("CLOUDFLARE_R2_ACCESS_ID")))
		fmt.Printf("R2 access key secret length: %d\n", len(os.Getenv("CLOUDFLARE_R2_ACCESS_KEY")))
		return fmt.Errorf("bucket does not exist or cannot be accessed: %v", err)
	}
	
	key = directory + "/" + key + ".jpg"
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/jpeg"),
	}

	// Log the bucket and key being used
	fmt.Printf("Uploading to R2 - Bucket: %s, Key: %s\n", s.bucket, key)
	
	_, err = s.s3Client.PutObject(ctx, input)
	if err != nil {
		fmt.Printf("R2 upload error: %v\n", err)
		return err
	}

	return nil
}

func (s *S3Service) FileExists(ctx context.Context, directory, slug string) (bool, error) {
	key := directory + "/" + slug + ".jpg"

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3Client.HeadObject(ctx, input)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *S3Service) DeleteFile(ctx context.Context, directory, slug string) error {
	key := directory + "/" + slug + ".jpg"
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3Client.DeleteObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// GetFileAWS GetFile GetBucket file from S3
func (s *S3Service) GetFileAWS(directory, slug string) (string, error) {
	key := directory + "/" + slug + ".jpg"
	return "https://pufacompsci.my.id/" + key, nil
}

func (s *S3Service) GetFileR2(directory, slug string) (string, error) {
	// Ensure we have a valid bucket name
	if s.bucket == "" {
		return "", fmt.Errorf("bucket name is empty")
	}
	
	// Format the URL correctly for R2
	// The URL format should be: https://sg.pufacomputing.live/{directory}/{slug}.jpg
	key := directory + "/" + slug + ".jpg"
	
	// Log the URL being generated
	url := "https://pufacompsci.my.id/" + key
	fmt.Printf("Generated R2 URL: %s\n", url)
	
	return url, nil
}
