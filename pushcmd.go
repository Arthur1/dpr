package dpr

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabriel-vasile/mimetype"
)

type PushCmd struct {
	File *os.File `arg:"" required:"" name:"path" help:"path of deployment package file"`
	Tags []string `name:"tag" short:"t" help:"tag for the deployment package"`
}

type Tag struct {
	Tag       string `dynamodbav:"tag"`
	ObjectKey string `dynamodbav:"object_key"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}

func (c *PushCmd) Run(globals *Globals) error {
	ctx := context.TODO()
	defer c.File.Close()

	cfg, err := globals.ReadConfig()
	if err != nil {
		return err
	}

	mimeType, err := mimetype.DetectReader(c.File)
	if err != nil {
		return err
	}
	if _, err := c.File.Seek(0, io.SeekStart); err != nil {
		return err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, c.File); err != nil {
		return err
	}
	hashTag := fmt.Sprintf("sha-%x", hash.Sum(nil))
	objectKey := fmt.Sprintf("%s%s", hashTag, filepath.Ext(c.File.Name()))

	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return err
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	dynamoDBClient := dynamodb.NewFromConfig(sdkConfig)

	if _, err := c.File.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if _, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(cfg.PackagesStore.S3BucketName),
		Key:         aws.String(objectKey),
		Body:        c.File,
		ContentType: aws.String(mimeType.String()),
	}); err != nil {
		return err
	}

	tags := slices.Concat([]string{hashTag}, c.Tags)
	slices.Sort(tags)
	uniqueTags := slices.Compact(tags)
	updatedAt := time.Now().Unix()
	for _, tag := range uniqueTags {
		tagRow := Tag{
			Tag:       tag,
			ObjectKey: objectKey,
			UpdatedAt: updatedAt,
		}
		tagItem, err := attributevalue.MarshalMap(tagRow)
		if err != nil {
			return err
		}

		if _, err := dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(cfg.TagsDB.DynamoDBTableName),
			Item:      tagItem,
		}); err != nil {
			return err
		}
	}
	return nil
}
