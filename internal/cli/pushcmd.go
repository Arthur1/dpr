package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/Arthur1/dpr/deploypackage"
	"github.com/Arthur1/dpr/internal/packagestore"
	"github.com/Arthur1/dpr/internal/tagdb"
	"github.com/Arthur1/dpr/usecase"
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

	packageStoreCli := packagestore.NewClientImpl(cfg.AwsConfig, cfg.PackageStore.S3BucketName)
	tagDBCli := tagdb.NewClientImpl(cfg.AwsConfig, cfg.TagDB.DynamoDBTableName)
	deploypackageRepository := deploypackage.NewDeployPackageRepositoryImpl(packageStoreCli, tagDBCli)
	pushUsecase := usecase.NewPushPackageUsecaseImpl(deploypackageRepository)

	result, err := pushUsecase.Push(ctx, &usecase.PushParam{
		Tags: c.Tags,
		File: c.File,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Pushed the deploy package. digest=%s\n", result.DeployPackage.Digest)
	return nil
}
