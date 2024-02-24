package usecase

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Arthur1/dpr/deploypackage"
	"github.com/Arthur1/dpr/lifecyclepolicy"
	"github.com/k1LoW/duration"
)

type ApplyLifecyclePolicyUsecase interface {
	ApplyLifecyclePolicy(ctx context.Context, param *ApplyLifecyclePolicyParam) (*ApplyLifecyclePolicyResult, error)
}

type ApplyLifecyclePolicyUsecaseImpl struct {
	deployPackageRepository deploypackage.DeployPackageRepository
}

var _ ApplyLifecyclePolicyUsecase = new(ApplyLifecyclePolicyUsecaseImpl)

func NewApplyLifecyclePolicyUsecaseImpl(
	deployPackageRepository deploypackage.DeployPackageRepository,
) *ApplyLifecyclePolicyUsecaseImpl {
	return &ApplyLifecyclePolicyUsecaseImpl{
		deployPackageRepository: deployPackageRepository,
	}
}

type ApplyLifecyclePolicyParam struct {
	DryRun          bool
	LifecyclePolicy *lifecyclepolicy.LifecyclePolicy
}

type ApplyLifecyclePolicyResult struct {
	ExpiredDeployPackages []*deploypackage.DeployPackage
}

func (u *ApplyLifecyclePolicyUsecaseImpl) ApplyLifecyclePolicy(ctx context.Context, param *ApplyLifecyclePolicyParam) (*ApplyLifecyclePolicyResult, error) {
	var targetDps []*deploypackage.DeployPackage
	for _, rule := range param.LifecyclePolicy.Rules {
		if rule.Action.Type != lifecyclepolicy.ActionTypeExpire {
			return nil, fmt.Errorf("\"%s\" is only supported for action.type value", lifecyclepolicy.ActionTypeExpire)
		}

		var (
			filterFunc func(int, int, *deploypackage.DeployPackage) bool
			err        error
		)
		switch rule.Selection.CountType {
		case lifecyclepolicy.CountTypeSincePackagePushed:
			filterFunc, err = buildFilterFuncSincePackagePushed(
				rule.Selection.CountUnit,
				rule.Selection.CountValue,
			)
			if err != nil {
				return nil, err
			}
		case lifecyclepolicy.CountTypePackageCountMoreThan:
			filterFunc, err = buildFilterFuncPackageCountMoreThan(
				rule.Selection.CountValue,
			)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("\"%s\" is not supported for selection.count-type value", rule.Selection.CountType)
		}

		var targetDpsByRule []*deploypackage.DeployPackage
		switch rule.Selection.TagStatus {
		case lifecyclepolicy.TagStatusUntagged:
			if targetDpsByRule, err = u.targetPackagesForUntaggedRule(ctx, filterFunc); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("\"%s\" is not supported for selection.tag-status value", rule.Selection.TagStatus)
		}
		targetDps = append(targetDps, targetDpsByRule...)
	}

	slices.SortFunc(targetDps, func(a, b *deploypackage.DeployPackage) int {
		return cmp.Compare(a.ObjectKey, b.ObjectKey)
	})
	targetDps = slices.CompactFunc(targetDps, func(a, b *deploypackage.DeployPackage) bool {
		return a.ObjectKey == b.ObjectBucket
	})

	result := &ApplyLifecyclePolicyResult{
		ExpiredDeployPackages: targetDps,
	}

	if param.DryRun || len(targetDps) == 0 {
		return result, nil
	}

	if err := u.deployPackageRepository.DeleteMultiple(ctx, targetDps); err != nil {
		return result, err
	}
	return result, nil
}

func (u *ApplyLifecyclePolicyUsecaseImpl) targetPackagesForUntaggedRule(
	ctx context.Context,
	filterFunc func(int, int, *deploypackage.DeployPackage) bool,
) ([]*deploypackage.DeployPackage, error) {
	dps, err := u.deployPackageRepository.GetUntaggedWithSortByUpdatedAtDesc(ctx)
	if err != nil {
		return nil, err
	}

	filteredDps := make([]*deploypackage.DeployPackage, 0, len(dps))
	for idx, dp := range dps {
		if filterFunc(len(dps), idx, dp) {
			filteredDps = append(filteredDps, dp)
		}
	}

	return filteredDps, nil
}

func buildFilterFuncSincePackagePushed(unit string, value int64) (func(int, int, *deploypackage.DeployPackage) bool, error) {
	dur, err := duration.Parse(fmt.Sprintf("%d%s", value, unit))
	if err != nil {
		return nil, err
	}
	now := time.Now()
	limit := now.Add(-dur)
	return func(_, _ int, dp *deploypackage.DeployPackage) bool {
		return time.Time(dp.UpdatedAt).Before(limit)
	}, nil
}

func buildFilterFuncPackageCountMoreThan(value int64) (func(int, int, *deploypackage.DeployPackage) bool, error) {
	return func(size, idx int, _ *deploypackage.DeployPackage) bool {
		selectSize := int64(size) - value
		if selectSize < 0 {
			selectSize = 0
		}
		return int64(idx) < selectSize
	}, nil
}
