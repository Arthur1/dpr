package dpr

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type TagStatus = string

const (
	TagStatusUntagged TagStatus = "untagged"
)

type CountType = string

const (
	CountTypeSincePackagePushed   CountType = "since-package-pushed"
	CountTypePackageCountMoreThan CountType = "package-count-more-than"
)

type ActionType = string

const (
	ActionTypeExpire ActionType = "expire"
)

type LifecyclePolicy struct {
	Rules []struct {
		Description string `yaml:"description" json:"description"`
		Selection   struct {
			TagStatus  string `yaml:"tag-status" json:"tag-status" jsonschema:"required,enum=untagged"`
			CountType  string `yaml:"count-type" json:"count-type" jsonschema:"required,enum=since-package-pushed,enum=package-count-more-than"`
			CountUnit  string `yaml:"count-unit" json:"count-unit"`
			CountValue int64  `yaml:"count-value" json:"count-value" jsonschema:"required"`
		} `yaml:"selection" json:"selection" jsonschema:"required"`
		Action struct {
			Type string `yaml:"type" json:"type" jsonschema:"required,enum=expire"`
		} `yaml:"action" json:"action" jsonschema:"required"`
	} `yaml:"rules" json:"rules" jsonschema:"required"`
}

func ReadLifecyclePolicy(file *os.File) (*LifecyclePolicy, error) {
	policy := &LifecyclePolicy{}
	b, _ := io.ReadAll(file)
	err := yaml.Unmarshal(b, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}
