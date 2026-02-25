package osutil

import (
	"os"
	"path/filepath"
	"syscall"
)

func Signals() []os.Signal {
	return []os.Signal{
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	}
}

func IsHangupSignal(signal os.Signal) bool {
	return signal == syscall.SIGHUP
}

func IsInterruptSignal(signal os.Signal) bool {
	return signal == os.Interrupt || signal == syscall.SIGINT || signal == syscall.SIGTERM
}

func WriteFile(file string, bytes []byte, perm os.FileMode) error {
	if err := CreateDir(filepath.Dir(file), 0755); err != nil {
		return err
	}
	return os.WriteFile(file, bytes, perm)
}

func CreateDir(dir string, perm os.FileMode) error {
	if !FileExists(dir) {
		if err := os.MkdirAll(dir, perm); err != nil {
			return err
		}
	}
	return nil
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}
