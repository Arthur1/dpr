package deploypackage

import (
	"io"
	"time"
)

type DeployPackage struct {
	Digest       string
	Tags         []string
	ObjectBucket string
	ObjectKey    string
	UpdatedAt    time.Time
	File         *DeployPackageFile
}

type DeployPackageFile struct {
	Body     io.Reader
	MimeType string
}

func (dp *DeployPackage) Copy() *DeployPackage {
	ndp := new(DeployPackage)
	*ndp = *dp
	return ndp
}
