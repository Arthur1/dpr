package dpr

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Region        string `yaml:"region"`
	PackagesStore struct {
		S3BucketName string `yaml:"s3-bucket-name"`
	} `yaml:"packages-store"`
	TagsDB struct {
		DynamoDBTableName string `yaml:"dynamodb-table-name"`
	} `yaml:"tags-db"`
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
