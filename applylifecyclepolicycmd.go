package dpr

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/Arthur1/dpr/tagdb"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/k1LoW/duration"
)

type ApplyLifecyclePolicyCmd struct {
	PolicyFile *os.File `arg:"" required:"" name:"policypath" help:"path of lifecycle policy file"`
	DryRun     bool     `name:"dry-run"`
}

func (c *ApplyLifecyclePolicyCmd) Run(globals *Globals) error {
	ctx := context.TODO()
	defer c.PolicyFile.Close()

	cfg, err := globals.ReadConfig()
	if err != nil {
		return err
	}

	policy, err := ReadLifecyclePolicy(c.PolicyFile)
	if err != nil {
		return err
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return err
	}
	tagDBClient := tagdb.NewClient(awsConfig, cfg.TagDB.DynamoDBTableName)

	var targetObjectKeys []string
	for _, rule := range policy.Rules {
		if rule.Action.Type != "expire" {
			return fmt.Errorf("\"expire\" is only supported for action.type value")
		}

		var filterFunc func(int, int, *tagdb.TagRow) bool
		switch rule.Selection.CountType {
		case "since-package-pushed":
			if filterFunc, err = buildFilterFuncSincePackagePushed(
				rule.Selection.CountUnit,
				rule.Selection.CountValue,
			); err != nil {
				return err
			}
		case "package-count-more-than":
			if filterFunc, err = buildFilterFuncPackageCountMoreThan(
				rule.Selection.CountValue,
			); err != nil {
				return err
			}
		default:
			return fmt.Errorf("\"%s\" is not supported for selection.count-type value", rule.Selection.CountType)
		}

		var targetObjectKeysPage []string
		switch rule.Selection.TagStatus {
		case "untagged":
			if targetObjectKeys, err = c.targetObjectsForUntaggedRule(ctx, tagDBClient, filterFunc); err != nil {
				return err
			}
		default:
			return fmt.Errorf("\"%s\" is not supported for selection.tag-status value", rule.Selection.TagStatus)
		}
		targetObjectKeys = append(targetObjectKeys, targetObjectKeysPage...)
	}

	slices.Sort(targetObjectKeys)
	targetObjectKeys = slices.Compact(targetObjectKeys)

	fmt.Printf("%+v\n", targetObjectKeys)
	if c.DryRun {
		return nil
	}

	return nil
}

func (c *ApplyLifecyclePolicyCmd) targetObjectsForUntaggedRule(
	ctx context.Context,
	tagDBClient *tagdb.Client,
	filterFunc func(int, int, *tagdb.TagRow) bool,
) ([]string, error) {
	tagRows, err := tagDBClient.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// it depends on the fact that string starting with @ comes first
	m := make(map[string]*tagdb.TagRow, len(tagRows))
	for _, tagRow := range tagRows {
		if strings.HasPrefix(tagRow.Tag, "@") {
			m[tagRow.ObjectKey] = tagRow
		} else {
			delete(m, tagRow.ObjectKey)
		}
	}

	tagRowsWithNoTags := make([]*tagdb.TagRow, 0, len(m))
	for _, tagRow := range m {
		tagRowsWithNoTags = append(tagRowsWithNoTags, tagRow)
	}
	slices.SortFunc(tagRowsWithNoTags, func(a, b *tagdb.TagRow) int {
		return cmp.Compare(time.Time(a.UpdatedAt).Unix(), time.Time(b.UpdatedAt).Unix())
	})

	filteredObjectKeys := make([]string, 0, len(tagRowsWithNoTags))
	for idx, tagRow := range tagRowsWithNoTags {
		if filterFunc(len(tagRowsWithNoTags), idx, tagRow) {
			filteredObjectKeys = append(filteredObjectKeys, tagRow.ObjectKey)
		}
	}

	return filteredObjectKeys, nil
}

func buildFilterFuncSincePackagePushed(unit string, value int64) (func(int, int, *tagdb.TagRow) bool, error) {
	dur, err := duration.Parse(fmt.Sprintf("%d%s", value, unit))
	if err != nil {
		return nil, err
	}
	now := time.Now()
	limit := now.Add(-dur)
	return func(_, _ int, tagRow *tagdb.TagRow) bool {
		return time.Time(tagRow.UpdatedAt).Before(limit)
	}, nil
}

func buildFilterFuncPackageCountMoreThan(value int64) (func(int, int, *tagdb.TagRow) bool, error) {
	return func(size, idx int, _ *tagdb.TagRow) bool {
		selectSize := int64(size) - value
		if selectSize < 0 {
			selectSize = 0
		}
		return int64(idx) < selectSize
	}, nil
}
