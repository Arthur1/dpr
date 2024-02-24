package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/Arthur1/dpr/deploypackage"
	"github.com/Arthur1/dpr/internal/packagestore"
	"github.com/Arthur1/dpr/internal/tagdb"
	"github.com/Arthur1/dpr/lifecyclepolicy"
	"github.com/Arthur1/dpr/usecase"
)

type ApplyLifecyclePolicyCmd struct {
	PolicyFile *os.File `arg:"" required:"" name:"policypath" help:"path of lifecycle policy file"`
	DryRun     bool     `name:"dry-run"`
}

func (c *ApplyLifecyclePolicyCmd) Run(globals *Globals) error {
	ctx := context.TODO()
	defer c.PolicyFile.Close()

	cfg, err := globals.ReadConfig(ctx)
	if err != nil {
		return err
	}

	packageStoreCli := packagestore.NewClientImpl(cfg.AwsConfig, cfg.PackageStore.S3BucketName)
	tagDBCli := tagdb.NewClientImpl(cfg.AwsConfig, cfg.TagDB.DynamoDBTableName)
	deploypackageRepository := deploypackage.NewDeployPackageRepositoryImpl(packageStoreCli, tagDBCli)
	applyUsecase := usecase.NewApplyLifecyclePolicyUsecaseImpl(deploypackageRepository)

	policy, err := lifecyclepolicy.ReadLifecyclePolicy(c.PolicyFile)
	if err != nil {
		return err
	}

	result, err := applyUsecase.ApplyLifecyclePolicy(ctx, &usecase.ApplyLifecyclePolicyParam{
		DryRun:          c.DryRun,
		LifecyclePolicy: policy,
	})
	if result != nil {
		if len(result.ExpiredDeployPackages) > 0 {
			fmt.Println("The following packages are expired:")
			for _, dp := range result.ExpiredDeployPackages {
				fmt.Printf("- %s\n", dp.ObjectKey)
			}
		} else {
			fmt.Println("No packages are expired.")
		}
	}
	if err != nil {
		return err
	}

	if !c.DryRun && len(result.ExpiredDeployPackages) > 0 {
		fmt.Println("Deleted.")
	}
	return nil
}
