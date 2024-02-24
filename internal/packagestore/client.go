package packagestore

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Client interface {
	Put(ctx context.Context, key string, contentType string, body io.Reader) error
	Find(ctx context.Context, key string) (*FindResult, error)
	Exists(ctx context.Context, key string) (bool, error)
	DeleteMultiple(ctx context.Context, keys []string) error
	GetBucketName(ctx context.Context) string
}

type ClientImpl struct {
	s3Cli      *s3.Client
	bucketName string
}

var _ Client = new(ClientImpl)

func NewClientImpl(cfg aws.Config, bucketName string) *ClientImpl {
	s3Cli := s3.NewFromConfig(cfg)
	return &ClientImpl{
		s3Cli:      s3Cli,
		bucketName: bucketName,
	}
}

func (c *ClientImpl) Put(ctx context.Context, key string, contentType string, body io.Reader) error {
	_, err := c.s3Cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

type FindResult struct {
	Body     io.ReadCloser
	MimeType string
}

func (c *ClientImpl) Find(ctx context.Context, key string) (*FindResult, error) {
	getObjectResult, err := c.s3Cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return &FindResult{
		Body:     getObjectResult.Body,
		MimeType: aws.ToString(getObjectResult.ContentType),
	}, nil
}

func (c *ClientImpl) Exists(ctx context.Context, key string) (bool, error) {
	if _, err := c.s3Cli.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (c *ClientImpl) DeleteMultiple(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	objects := make([]types.ObjectIdentifier, 0, len(keys))
	for _, key := range keys {
		objects = append(objects, types.ObjectIdentifier{Key: aws.String(key)})
	}

	_, err := c.s3Cli.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(c.bucketName),
		Delete: &types.Delete{Objects: objects},
	})
	return err
}

func (c *ClientImpl) GetBucketName(ctx context.Context) string {
	return c.bucketName
}
