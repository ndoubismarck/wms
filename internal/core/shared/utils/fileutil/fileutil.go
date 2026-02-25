package fileutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}

func CopyFile(source, dest string) error {
	if !Exists(source) {
		return fmt.Errorf("source file does not exist: %s", source)
	}
	destDir := filepath.Dir(source)
	if !Exists(source) {
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return err
		}
	}
	in, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)
	if _, err = io.Copy(out, bytes.NewReader(in)); err != nil {
		return err
	}
	return nil
}

func Delete(path string) error {
	return os.Remove(path)
}
