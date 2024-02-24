package deploypackage

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Arthur1/dpr/internal/packagestore"
	"github.com/Arthur1/dpr/internal/tagdb"
)

type DeployPackageRepository interface {
	Save(ctx context.Context, dp *DeployPackage) error
	FindByTag(ctx context.Context, tag string) (*DeployPackage, error)
	FindByDigest(ctx context.Context, digest string) (*DeployPackage, error)
	GetUntaggedWithSortByUpdatedAtDesc(ctx context.Context) ([]*DeployPackage, error)
	LoadFile(ctx context.Context, dp *DeployPackage) (*DeployPackage, error)
	DeleteMultiple(ctx context.Context, dps []*DeployPackage) error
}

type DeployPackageRepositoryImpl struct {
	packageStoreCli packagestore.Client
	tagDBCli        tagdb.Client
}

var _ DeployPackageRepository = new(DeployPackageRepositoryImpl)

func NewDeployPackageRepositoryImpl(
	packageStoreCli packagestore.Client,
	tagDBCli tagdb.Client,
) *DeployPackageRepositoryImpl {
	return &DeployPackageRepositoryImpl{
		packageStoreCli: packageStoreCli,
		tagDBCli:        tagDBCli,
	}
}

func (r *DeployPackageRepositoryImpl) Save(ctx context.Context, dp *DeployPackage) error {
	if err := r.packageStoreCli.Put(ctx, dp.ObjectKey, dp.File.MimeType, dp.File.Body); err != nil {
		return err
	}
	tagRows := make([]*tagdb.TagRow, 0, len(dp.Tags)+1)
	tagRows = append(tagRows, &tagdb.TagRow{
		Type:      "digest",
		Tag:       fmt.Sprintf("@%s", dp.Digest),
		ObjectKey: dp.ObjectKey,
		UpdatedAt: dp.UpdatedAt,
	})
	for _, tag := range dp.Tags {
		tagRows = append(tagRows, &tagdb.TagRow{
			Type:      "tag",
			Tag:       tag,
			ObjectKey: dp.ObjectKey,
			UpdatedAt: dp.UpdatedAt,
		})
	}
	if err := r.tagDBCli.PutMultiple(ctx, tagRows); err != nil {
		return err
	}
	return nil
}

func (r *DeployPackageRepositoryImpl) FindByTag(ctx context.Context, tag string) (*DeployPackage, error) {
	tagRow, err := r.tagDBCli.FindByTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	_, err = r.packageStoreCli.Exists(ctx, tagRow.ObjectKey)
	if err != nil {
		return nil, err
	}
	tagRows, err := r.tagDBCli.GetByObjectKey(ctx, tagRow.ObjectKey)
	if err != nil {
		return nil, err
	}
	digest, tags := convertTagRowsToDigestAndTag(tagRows)

	lastUpdatedAt := getLastUpdatedAtFromTagRows(
		slices.Concat(tagRows, []*tagdb.TagRow{tagRow}),
	)

	return &DeployPackage{
		Digest:       digest,
		Tags:         tags,
		ObjectBucket: r.packageStoreCli.GetBucketName(ctx),
		ObjectKey:    tagRow.ObjectKey,
		UpdatedAt:    lastUpdatedAt,
	}, nil
}

func (r *DeployPackageRepositoryImpl) FindByDigest(ctx context.Context, digest string) (*DeployPackage, error) {
	tagRow, err := r.tagDBCli.FindByDigest(ctx, digest)
	if err != nil {
		return nil, err
	}
	_, err = r.packageStoreCli.Exists(ctx, tagRow.ObjectKey)
	if err != nil {
		return nil, err
	}
	tagRows, err := r.tagDBCli.GetByObjectKey(ctx, tagRow.ObjectKey)
	if err != nil {
		return nil, err
	}
	digest, tags := convertTagRowsToDigestAndTag(tagRows)

	lastUpdatedAt := getLastUpdatedAtFromTagRows(
		slices.Concat(tagRows, []*tagdb.TagRow{tagRow}),
	)
	return &DeployPackage{
		Digest:       digest,
		Tags:         tags,
		ObjectBucket: r.packageStoreCli.GetBucketName(ctx),
		ObjectKey:    tagRow.ObjectKey,
		UpdatedAt:    lastUpdatedAt,
	}, nil
}

func (r *DeployPackageRepositoryImpl) GetUntaggedWithSortByUpdatedAtDesc(ctx context.Context) ([]*DeployPackage, error) {
	tagRows, err := r.tagDBCli.GetAll(ctx)
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

	dps := make([]*DeployPackage, 0, len(m))
	for _, tagRow := range m {
		dps = append(dps, &DeployPackage{
			Digest:       strings.Replace(tagRow.Tag, "@", "", 1),
			Tags:         []string{},
			ObjectBucket: r.packageStoreCli.GetBucketName(ctx),
			ObjectKey:    tagRow.ObjectKey,
			UpdatedAt:    tagRow.UpdatedAt,
		})
	}

	slices.SortFunc(dps, func(a, b *DeployPackage) int {
		return cmp.Compare(a.UpdatedAt.Unix(), b.UpdatedAt.Unix())
	})

	return dps, nil
}

func (r *DeployPackageRepositoryImpl) LoadFile(ctx context.Context, dp *DeployPackage) (*DeployPackage, error) {
	result, err := r.packageStoreCli.Find(ctx, dp.ObjectKey)
	if err != nil {
		return nil, err
	}
	ndp := dp.Copy()
	ndp.File = &DeployPackageFile{
		Body:     result.Body,
		MimeType: result.MimeType,
	}
	return ndp, nil
}

func convertTagRowsToDigestAndTag(tagRows []*tagdb.TagRow) (string, []string) {
	digest := ""
	tags := make([]string, 0, len(tagRows))
	for _, row := range tagRows {
		switch row.Type {
		case "digest":
			digest = strings.Replace(row.Tag, "@", "", 1)
		case "tag":
			tags = append(tags, row.Tag)
		}
	}
	return digest, tags
}

func getLastUpdatedAtFromTagRows(tagRows []*tagdb.TagRow) time.Time {
	updatedAts := make([]time.Time, 0, len(tagRows))
	for _, tagRow := range tagRows {
		updatedAts = append(updatedAts, tagRow.UpdatedAt)
	}
	lastUpdateAt := slices.MaxFunc(updatedAts, func(a, b time.Time) int {
		return cmp.Compare(a.Unix(), b.Unix())
	})
	return lastUpdateAt
}

func (r *DeployPackageRepositoryImpl) DeleteMultiple(ctx context.Context, dps []*DeployPackage) error {
	objectKeys := make([]string, 0, len(dps))
	for _, dp := range dps {
		objectKeys = append(objectKeys, dp.ObjectKey)
	}
	if err := r.packageStoreCli.DeleteMultiple(ctx, objectKeys); err != nil {
		return err
	}
	tagRows := []*tagdb.TagRow{}
	for _, dp := range dps {
		tagRows = append(tagRows, &tagdb.TagRow{
			Type:      "digest",
			Tag:       fmt.Sprintf("@%s", dp.Digest),
			ObjectKey: dp.ObjectKey,
			UpdatedAt: dp.UpdatedAt,
		})
		for _, tag := range dp.Tags {
			tagRows = append(tagRows, &tagdb.TagRow{
				Type:      "tag",
				Tag:       tag,
				ObjectKey: dp.ObjectKey,
				UpdatedAt: dp.UpdatedAt,
			})
		}
	}
	if err := r.tagDBCli.DeleteMultiple(ctx, tagRows); err != nil {
		return err
	}
	return nil
}
