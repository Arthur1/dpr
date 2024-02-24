package packagefile

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

type PackageFile struct {
	DigestType  string
	DigestValue string
	Body        io.Reader
	Ext         string
	MimeType    string
}

func NewPackageFile(file *os.File) (*PackageFile, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}
	hashSum := hash.Sum(nil)
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return &PackageFile{
		DigestType:  "sha256",
		DigestValue: fmt.Sprintf("%x", hashSum),
		Body:        file,
		Ext:         filepath.Ext(file.Name()),
		MimeType:    mimeType.String(),
	}, nil
}

type StoredPackageFile struct {
	Body io.ReadCloser
}
