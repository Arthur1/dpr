package usecase

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/Arthur1/dpr/deploypackage"
	"github.com/Arthur1/dpr/packagefile"
)

type PushPackageUsecase interface {
	Push(ctx context.Context, param *PushParam) (*PushResult, error)
}

type PushPackageUsecaseImpl struct {
	deployPackageRepository deploypackage.DeployPackageRepository
}

var _ PushPackageUsecase = new(PushPackageUsecaseImpl)

func NewPushPackageUsecaseImpl(deployPackageRepository deploypackage.DeployPackageRepository) *PushPackageUsecaseImpl {
	return &PushPackageUsecaseImpl{
		deployPackageRepository: deployPackageRepository,
	}
}

type PushParam struct {
	Tags []string
	File *os.File
}

type PushResult struct {
	DeployPackage *deploypackage.DeployPackage
}

func (u *PushPackageUsecaseImpl) Push(ctx context.Context, param *PushParam) (*PushResult, error) {
	file, err := packagefile.NewPackageFile(param.File)
	if err != nil {
		return nil, err
	}

	slices.Sort(param.Tags)
	uniqueTags := slices.Compact(param.Tags)

	dp := &deploypackage.DeployPackage{
		Digest:       fmt.Sprintf("%s:%s", file.DigestType, file.DigestValue),
		Tags:         uniqueTags,
		ObjectBucket: "", // don't care
		ObjectKey:    fmt.Sprintf("%s-%s%s", file.DigestType, file.DigestValue, file.Ext),
		File: &deploypackage.DeployPackageFile{
			Body:     file.Body,
			MimeType: file.MimeType,
		},
	}
	u.deployPackageRepository.Save(ctx, dp)
	return &PushResult{
		DeployPackage: dp,
	}, nil
}
