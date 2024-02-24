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

	"github.com/Arthur1/dpr/packagestore"
	"github.com/Arthur1/dpr/tagdb"
	"github.com/aws/aws-sdk-go-v2/config"
)

type PushCmd struct {
	File *os.File `arg:"" required:"" name:"path" help:"path of deployment package file"`
	Tags []string `name:"tag" short:"t" help:"tag for the deployment package"`
}

func (c *PushCmd) Run(globals *Globals) error {
	ctx := context.TODO()
	defer c.File.Close()

	cfg, err := globals.ReadConfig(ctx)
	if err != nil {
		return err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, c.File); err != nil {
		return err
	}
	hashSum := hash.Sum(nil)
	hashTag := fmt.Sprintf("@sha256:%x", hash.Sum(nil))
	objectKey := fmt.Sprintf("sha256-%x%s", hashSum, filepath.Ext(c.File.Name()))

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return err
	}
	packageStoreClient := packagestore.NewClient(awsConfig, cfg.PackageStore.S3BucketName)
	tagDBClient := tagdb.NewClient(awsConfig, cfg.TagDB.DynamoDBTableName)

	if _, err := c.File.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := packageStoreClient.PutPackage(ctx, &packagestore.PutPackageInput{
		File:      c.File,
		ObjectKey: objectKey,
	}); err != nil {
		return err
	}

	tags := slices.Concat([]string{hashTag}, c.Tags)
	slices.Sort(tags)
	uniqueTags := slices.Compact(tags)
	updatedAt := time.Now()

	if err := tagDBClient.PutTags(ctx, &tagdb.PutTagsInput{
		Tags:      uniqueTags,
		ObjectKey: objectKey,
		UpdatedAt: updatedAt,
	}); err != nil {
		return err
	}
	return nil
}
