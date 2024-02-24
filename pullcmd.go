package dpr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Arthur1/dpr/packagestore"
	"github.com/Arthur1/dpr/tagdb"
	"github.com/aws/aws-sdk-go-v2/config"
)

type PullCmd struct {
	Tag       string `name:"tag" short:"t" help:"tag for search condition"`
	Digest    string `name:"digest" short:"d" help:"digest for search condition"`
	CheckOnly bool   `name:"check-only" help:"checks for the existence of the package file, but does not download it"`
}

type PullCmdOutput struct {
	S3Bucket string `json:"s3_bucket"`
	S3Key    string `json:"s3_key"`
}

func (c *PullCmd) Run(globals *Globals) error {
	ctx := context.TODO()

	cfg, err := globals.ReadConfig(ctx)
	if err != nil {
		return err
	}

	tag := c.Tag
	if c.Digest != "" {
		tag = fmt.Sprintf("@%s", c.Digest)
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return err
	}
	packageStoreClient := packagestore.NewClient(awsConfig, cfg.PackageStore.S3BucketName)
	tagDBClient := tagdb.NewClient(awsConfig, cfg.TagDB.DynamoDBTableName)

	tagRow, err := tagDBClient.FindByTag(ctx, &tagdb.FindByTagInput{
		Tag: tag,
	})
	if err != nil {
		return err
	}
	objectKey := tagRow.ObjectKey

	if c.CheckOnly {
		if _, err := packageStoreClient.ExistsPackage(ctx, &packagestore.ExistsPackageInput{
			ObjectKey: objectKey,
		}); err != nil {
			return err
		}
	} else {
		body, err := packageStoreClient.FindPackage(ctx, &packagestore.FindPackageInput{
			ObjectKey: objectKey,
		})
		if err != nil {
			return err
		}
		defer body.Close()

		_, filename := filepath.Split(objectKey)
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(file, body); err != nil {
			return err
		}
	}

	output := PullCmdOutput{
		S3Bucket: cfg.PackageStore.S3BucketName,
		S3Key:    objectKey,
	}
	outputJson, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Println(string(outputJson))

	return nil
}
