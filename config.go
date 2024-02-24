package dpr

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/yaml.v3"
)

type Config struct {
	*RawConfig
	AwsConfig aws.Config
}

type RawConfig struct {
	Region       string `yaml:"region" json:"region" jsonschema:"required,description=AWS region of dpr resources,example=ap-northeast-1"`
	PackageStore struct {
		S3BucketName string `yaml:"s3-bucket-name" json:"s3-bucket-name" jsonschema:"required,description=Amazon S3 bucket name for dpr packages store"`
	} `yaml:"package-store" json:"package-store"`
	TagDB struct {
		DynamoDBTableName string `yaml:"dynamodb-table-name" json:"dynamodb-table-name" jsonschema:"required,description=Amazon DynamoDB table name for dpr tags database"`
	} `yaml:"tag-db" json:"tag-db"`
}

func (g *Globals) ReadConfig(ctx context.Context) (*Config, error) {
	rawcfg := &RawConfig{}
	b, _ := os.ReadFile(g.ConfigFile)
	err := yaml.Unmarshal(b, rawcfg)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		RawConfig: rawcfg,
	}
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}
	cfg.AwsConfig = awsConfig
	return cfg, nil
}
