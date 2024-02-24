package packagestore

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
)

type Client struct {
	s3Cli      *s3.Client
	bucketName string
}

func NewClient(cfg aws.Config, bucketName string) *Client {
	s3Cli := s3.NewFromConfig(cfg)
	return &Client{
		s3Cli:      s3Cli,
		bucketName: bucketName,
	}
}

type PutPackageInput struct {
	File      *os.File
	ObjectKey string
}

func (c *Client) PutPackage(ctx context.Context, input *PutPackageInput) error {
	mimeType, err := mimetype.DetectReader(input.File)
	if err != nil {
		return err
	}
	if _, err := input.File.Seek(0, io.SeekStart); err != nil {
		return err
	}

	_, err = c.s3Cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(input.ObjectKey),
		Body:        input.File,
		ContentType: aws.String(mimeType.String()),
	})
	return err
}

type FindPackageInput struct {
	ObjectKey string
}

func (c *Client) FindPackage(ctx context.Context, input *FindPackageInput) (io.ReadCloser, error) {
	getObjectResult, err := c.s3Cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(input.ObjectKey),
	})
	if err != nil {
		return nil, err
	}
	return getObjectResult.Body, nil
}

type ExistsPackageInput struct {
	ObjectKey string
}

func (c *Client) ExistsPackage(ctx context.Context, input *ExistsPackageInput) (bool, error) {
	if _, err := c.s3Cli.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(input.ObjectKey),
	}); err != nil {
		return false, err
	}
	return true, nil
}

type DeletePackagesInput struct {
	ObjectKeys []string
}

func (c *Client) DeletePackages(ctx context.Context, input *DeletePackagesInput) error {
	if len(input.ObjectKeys) == 0 {
		return nil
	}

	objectIds := make([]types.ObjectIdentifier, 0, len(input.ObjectKeys))
	for _, objectKey := range input.ObjectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(objectKey)})
	}

	_, err := c.s3Cli.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(c.bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	return err
}
