package usecase

import (
	"context"

	"github.com/Arthur1/dpr/deploypackage"
)

type PullPackageUsecase interface {
	PullByDigest(ctx context.Context, param *PullByDigestParam) (*PullResult, error)
	PullByTag(ctx context.Context, param *PullByTagParam) (*PullResult, error)
}

type PullPackageUsecaseImpl struct {
	deployPackageRepository deploypackage.DeployPackageRepository
}

var _ PullPackageUsecase = new(PullPackageUsecaseImpl)

func NewPullPackageUsecaseImpl(
	deployPackageRepository deploypackage.DeployPackageRepository,
) *PullPackageUsecaseImpl {
	return &PullPackageUsecaseImpl{
		deployPackageRepository: deployPackageRepository,
	}
}

type PullByDigestParam struct {
	Digest    string
	NeedsFile bool
}

type PullResult struct {
	DeployPackage *deploypackage.DeployPackage
}

func (u *PullPackageUsecaseImpl) PullByDigest(ctx context.Context, param *PullByDigestParam) (*PullResult, error) {
	dp, err := u.deployPackageRepository.FindByDigest(ctx, param.Digest)
	if err != nil {
		return nil, err
	}
	if param.NeedsFile {
		dp, err = u.deployPackageRepository.LoadFile(ctx, dp)
		if err != nil {
			return nil, err
		}
	}
	return &PullResult{
		DeployPackage: dp,
	}, nil
}

type PullByTagParam struct {
	Tag       string
	NeedsFile bool
}

func (u *PullPackageUsecaseImpl) PullByTag(ctx context.Context, param *PullByTagParam) (*PullResult, error) {
	dp, err := u.deployPackageRepository.FindByTag(ctx, param.Tag)
	if err != nil {
		return nil, err
	}
	if param.NeedsFile {
		dp, err = u.deployPackageRepository.LoadFile(ctx, dp)
		if err != nil {
			return nil, err
		}
	}
	return &PullResult{
		DeployPackage: dp,
	}, nil
}
