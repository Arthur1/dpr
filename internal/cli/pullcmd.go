package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Arthur1/dpr/deploypackage"
	"github.com/Arthur1/dpr/internal/packagestore"
	"github.com/Arthur1/dpr/internal/tagdb"
	"github.com/Arthur1/dpr/usecase"
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

	packageStoreCli := packagestore.NewClientImpl(cfg.AwsConfig, cfg.PackageStore.S3BucketName)
	tagDBCli := tagdb.NewClientImpl(cfg.AwsConfig, cfg.TagDB.DynamoDBTableName)
	deploypackageRepository := deploypackage.NewDeployPackageRepositoryImpl(packageStoreCli, tagDBCli)
	pullUsecase := usecase.NewPullPackageUsecaseImpl(deploypackageRepository)

	var result *usecase.PullResult
	if c.Tag != "" {
		result, err = pullUsecase.PullByTag(ctx, &usecase.PullByTagParam{
			Tag:       c.Tag,
			NeedsFile: !c.CheckOnly,
		})
	} else {
		result, err = pullUsecase.PullByDigest(ctx, &usecase.PullByDigestParam{
			Digest:    c.Digest,
			NeedsFile: !c.CheckOnly,
		})
	}
	if err != nil {
		return err
	}

	output := PullCmdOutput{
		S3Bucket: result.DeployPackage.ObjectBucket,
		S3Key:    result.DeployPackage.ObjectKey,
	}
	outputJson, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Println(string(outputJson))

	if result.DeployPackage.File != nil {
		_, filename := filepath.Split(result.DeployPackage.ObjectKey)
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(file, result.DeployPackage.File.Body); err != nil {
			return err
		}
	}

	return nil
}
