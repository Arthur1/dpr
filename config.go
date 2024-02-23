package dpr

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Region        string `yaml:"region" json:"region" jsonschema:"required,description=AWS region of dpr resources,example=ap-northeast-1"`
	PackagesStore struct {
		S3BucketName string `yaml:"s3-bucket-name" json:"s3-bucket-name" jsonschema:"required,description=Amazon S3 bucket name for dpr packages store"`
	} `yaml:"packages-store" json:"packages-store"`
	TagsDB struct {
		DynamoDBTableName string `yaml:"dynamodb-table-name" json:"dynamodb-table-name" jsonschema:"required,description=Amazon DynamoDB table name for dpr tags database"`
	} `yaml:"tags-db" json:"tags-db"`
}

func (g *Globals) ReadConfig() (*Config, error) {
	config := &Config{}
	b, _ := os.ReadFile(g.ConfigFile)
	err := yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
