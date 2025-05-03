package services

import (
	"Backend/configs"
	"bytes"
	"context"
	"fmt"

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
	var url = s3Config.S3Endpoint
	var usePathStyle = s3Config.S3UsePathStyle

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           url,
					SigningRegion: region,
				}, nil
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create an Amazon S3 service client access key and so on
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = usePathStyle
	})

	return &S3Service{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func NewR2Service() (*S3Service, error) {
	s3Config := configs.LoadConfig()
	var bucket = s3Config.S3Bucket
	// Set a default bucket name if it's empty
	if bucket == "" {
		bucket = "pufa-2025"
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

	// Use path style addressing to avoid bucket name validation issues
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

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
	// Format the key as directory/key.jpg
	key = directory + "/" + key + ".jpg"
	
	// Log the bucket and key for debugging
	fmt.Printf("Uploading to R2 - Bucket: %s, Key: %s\n", s.bucket, key)
	
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/jpeg"),
		// Public read access is handled by bucket policy in Cloudflare R2
	}

	_, err := s.s3Client.PutObject(ctx, input)
	if err != nil {
		fmt.Printf("Error uploading to R2: %v\n", err)
		return err
	}

	fmt.Printf("Successfully uploaded to R2 - Bucket: %s, Key: %s\n", s.bucket, key)
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
	return "https://id.pufacomputing.live/" + key, nil
}

func (s *S3Service) GetFileR2(directory, slug string) (string, error) {
	// Format the key as it's stored in R2
	key := directory + "/" + slug + ".jpg"
	
	// For Cloudflare R2 with custom domain
	// Use the public URL format that works with your Cloudflare R2 setup
	return fmt.Sprintf("https://pufacompsci.my.id/%s", key), nil
}
