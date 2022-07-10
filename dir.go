package tutil

import (
	"fmt"
	"os"
	"path/filepath"
)

var TempDirName = "tutil"

type DirOption func(root string) error

func WithSubDirs(subDirCount int) DirOption {
	return func(root string) error {
		if subDirCount < 1 {
			return fmt.Errorf("WithSubDirs() expects a positive subDirCount, got %q", subDirCount)
		}

		subDirs := []string{root}
		for i := 0; i < subDirCount; i++ {
			subDirs = append(subDirs, fmt.Sprintf("level-%d", i+1))
		}

		dirName := filepath.Join(subDirs...)
		if err := os.MkdirAll(dirName, 0770); err != nil {
			return err
		}

		return nil
	}
}

func TempDir(options ...DirOption) (string, func() error, error) {
	dir, err := os.MkdirTemp("", TempDirName)
	if err != nil {
		return "", nil, fmt.Errorf("couldn't set up temp dir: %w", err)
	}

	for _, option := range options {
		err := option(dir)
		if err != nil {
			return "", nil, fmt.Errorf("couldn't apply option during temp dir creation: %w", err)
		}
	}

	cleanup := func() error {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("couldn't remove temp dir during cleanup: %w", err)
		}
		return nil
	}

	return dir, cleanup, nil
}
