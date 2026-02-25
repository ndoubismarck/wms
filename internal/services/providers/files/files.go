package files

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/fileutil"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

type Provider struct {
	ctx types.IContext
}

func New(ctx types.IContext) *Provider {
	return &Provider{
		ctx: ctx,
	}
}

func (p *Provider) Cleanup() error {
	dir := UploadsDir
	if !fileutil.Exists(dir) {
		return nil
	}
	now := time.Now().UTC()
	maxAge := time.Hour * 24 * 7
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			path := filepath.Join(dir, entry.Name())
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("failed to get file info for %s: %v", path, err)
			}
			if now.Sub(info.ModTime().UTC()) > maxAge {
				if err = os.Remove(path); err != nil {
					return fmt.Errorf("failed to delete file %s: %v", path, err)
				}
			}
		}
	}
	return nil
}

func (p *Provider) upload(data UploadFileData, ext string) (*UploadFileResult, error) {
	dest := UploadsDir
	if !fileutil.Exists(dest) {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return nil, err
		}
	}
	name := fmt.Sprintf("%s%s", ksuid.New().String(), ext)
	dest = filepath.Join(dest, name)
	file, err := data.File.Open()
	if err != nil {
		return nil, err
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	out, err := os.Create(dest)
	if err != nil {
		return nil, err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)
	if _, err = io.Copy(out, file); err != nil {
		return nil, err
	}
	return &UploadFileResult{
		Size: data.File.Size,
		Type: data.File.Header.Get("Content-Type"),
		Path: dest,
		Name: data.File.Filename,
	}, nil
}

func (p *Provider) Upload(data UploadFileData) (*UploadFileResult, error) {
	var ext string
	contentType := data.File.Header.Get("Content-Type")
	switch contentType {
	case "image/png":
		ext = ".png"
	case "image/jpg":
		ext = ".jpg"
	case "image/jpeg":
		ext = ".jpg"
	}
	if len(ext) == 0 {
		return nil, errors.New("unsupported image format")
	}
	return p.upload(data, ext)
}

func (p *Provider) DeleteFile(path string) error {
	if fileutil.Exists(path) {
		return fileutil.Delete(path)
	}
	return nil
}

func (p *Provider) DeleteFiles(paths []string) error {
	for _, val := range paths {
		if err := p.DeleteFile(val); err != nil {
			return nil
		}
	}
	return nil
}

func (p *Provider) CopyUploadedFile(path string, dest string) (string, error) {
	if !fileutil.Exists(dest) {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return "", err
		}
	}
	ext := filepath.Ext(filepath.Base(path))
	if !strings.HasPrefix(ext, ".") {
		ext = fmt.Sprintf(".%s", ext)
	}
	name := fmt.Sprintf("%s%s", cryptoutil.UUID(), filepath.Base(path))
	name = fmt.Sprintf("%s%s", cryptoutil.Sha1(name), ext)
	dest = filepath.Join(dest, name)
	if err := fileutil.CopyFile(path, dest); err != nil {
		return "", err
	}
	return dest, nil
}

func (p *Provider) CopyUploadedFiles(paths []string, dest string) ([]string, error) {
	var moved []string
	for _, val := range paths {
		path, err := p.CopyUploadedFile(val, dest)
		if err != nil {
			return nil, nil
		}
		moved = append(moved, path)
	}
	return moved, nil
}
