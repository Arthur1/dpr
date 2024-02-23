package dpr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PullCmd struct {
	Tag       string `name:"tag" short:"t" required:"" help:"tag for search condition"`
	CheckOnly bool   `name:"check-only" help:"checks for the existence of the package file, but does not download it"`
}

type PullCmdOutput struct {
	S3Bucket string `json:"s3_bucket"`
	S3Key    string `json:"s3_key"`
}

func (c *PullCmd) Run(globals *Globals) error {
	ctx := context.TODO()

	cfg, err := globals.ReadConfig()
	if err != nil {
		return err
	}

	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return err
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	dynamoDBClient := dynamodb.NewFromConfig(sdkConfig)

	getItemResult, err := dynamoDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(cfg.TagsDB.DynamoDBTableName),
		Key: map[string]types.AttributeValue{
			"tag": &types.AttributeValueMemberS{
				Value: c.Tag,
			},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return err
	}
	tagRow := Tag{}
	if err = attributevalue.UnmarshalMap(getItemResult.Item, &tagRow); err != nil {
		return err
	}

	if c.CheckOnly {
		if _, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(cfg.PackagesStore.S3BucketName),
			Key:    aws.String(tagRow.ObjectKey),
		}); err != nil {
			return err
		}
	} else {
		getObjectResult, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(cfg.PackagesStore.S3BucketName),
			Key:    aws.String(tagRow.ObjectKey),
		})
		if err != nil {
			return err
		}
		defer getObjectResult.Body.Close()

		_, filename := filepath.Split(tagRow.ObjectKey)
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(file, getObjectResult.Body); err != nil {
			return err
		}
	}

	output := PullCmdOutput{
		S3Bucket: cfg.PackagesStore.S3BucketName,
		S3Key:    tagRow.ObjectKey,
	}
	outputJson, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Println(string(outputJson))

	return nil
}
