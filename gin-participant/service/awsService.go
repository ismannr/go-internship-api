package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct {
	PresignClient *s3.PresignClient
}

func NewPresigner(config aws.Config) *Presigner {
	return &Presigner{
		PresignClient: s3.NewPresignClient(s3.NewFromConfig(config)),
	}
}

func getConfig() (aws.Config, error) {
	awsAccessKeyID := os.Getenv("AWS_S3_ACCESS_KEY")
	awsSecretAccessKey := os.Getenv("AWS_S3_SECRET_KEY")
	awsRegion := os.Getenv("AWS_S3_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     awsAccessKeyID,
				SecretAccessKey: awsSecretAccessKey,
			}, nil
		})),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

func (presigner *Presigner) GetObject(bucketName, objectKey string, lifetimeSecs int64) (string, error) {
	req, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(int64(lifetimeSecs) * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
		return "", err
	}
	return req.URL, nil
}

func (presigner *Presigner) PutObject(bucketName, objectKey string, lifetimeSecs int64) (string, error) {
	req, err := presigner.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(int64(lifetimeSecs) * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
		return "", err
	}
	return req.URL, nil
}

func (presigner *Presigner) DeleteObject(bucketName, objectKey string, lifetimeSecs int64) (string, error) {
	req, err := presigner.PresignClient.PresignDeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(int64(lifetimeSecs) * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to delete object %v. Here's why: %v\n", objectKey, err)
		return "", err
	}
	return req.URL, nil
}
