package dpr

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type LifecyclePolicy struct {
	Rules []struct {
		Description string `yaml:"description"`
		Selection   struct {
			TagStatus  string `yaml:"tag-status"`
			CountType  string `yaml:"count-type"`
			CountUnit  string `yaml:"count-unit"`
			CountValue int64  `yaml:"count-value"`
		} `yaml:"selection"`
		Action struct {
			Type string `yaml:"type"`
		}
	} `yaml:"rules"`
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
