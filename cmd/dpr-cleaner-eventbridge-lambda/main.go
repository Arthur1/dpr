package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/Arthur1/dpr/deploypackage"
	"github.com/Arthur1/dpr/internal/packagestore"
	"github.com/Arthur1/dpr/internal/tagdb"
	"github.com/Arthur1/dpr/lifecyclepolicy"
	"github.com/Arthur1/dpr/usecase"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	lambda.Start(handler)
}

type Event struct {
	LifecyclePolicy lifecyclepolicy.LifecyclePolicy `json:"lifecycle-policy"`
	PackageStore    struct {
		S3BucketName string `json:"s3-bucket-name"`
	} `json:"package-store"`
	TagDB struct {
		DynamoDBTableName string `json:"dynamodb-table-name"`
	} `json:"tag-db"`
}

func handler(ctx context.Context, event Event) (string, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	region := os.Getenv("AWS_REGION")
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", err
	}

	packageStoreCli := packagestore.NewClientImpl(awsConfig, event.PackageStore.S3BucketName)
	tagDBCli := tagdb.NewClientImpl(awsConfig, event.TagDB.DynamoDBTableName)
	deploypackageRepository := deploypackage.NewDeployPackageRepositoryImpl(packageStoreCli, tagDBCli)
	applyUsecase := usecase.NewApplyLifecyclePolicyUsecaseImpl(deploypackageRepository)

	result, err := applyUsecase.ApplyLifecyclePolicy(ctx, &usecase.ApplyLifecyclePolicyParam{
		LifecyclePolicy: &event.LifecyclePolicy,
	})
	if err != nil {
		return "", err
	}
	digests := make([]string, 0, len(result.ExpiredDeployPackages))
	for _, dp := range result.ExpiredDeployPackages {
		digests = append(digests, dp.Digest)
	}
	slog.Info(
		"success",
		slog.Int("expiredPackagesNumber", len(result.ExpiredDeployPackages)),
		slog.String("expiredPackageDigests", strings.Join(digests, ",")),
	)
	return "", nil
}
