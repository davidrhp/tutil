package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTempDir(t *testing.T) {
	type args struct {
		options []DirOption
	}
	tests := []struct {
		name    string
		args    args
		wantDir string
		wantErr bool
	}{
		{
			name: "create temp dir and clean up afterwards",
			args: args{
				options: nil,
			},
			wantDir: TempDirName,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDir, gotCleanup, err := TempDir(tt.args.options...)

			// check error
			if err != nil != tt.wantErr {
				t.Fatalf("TempDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			// check correct directory name
			if !strings.Contains(filepath.Base(gotDir), tt.wantDir) {
				t.Errorf("TempDir() gotDir = %v, wantDir %v", gotDir, tt.wantDir)
			}
			// check temp dir in default temp dir location.
			if !strings.Contains(os.TempDir(), filepath.Dir(gotDir)) {
				t.Errorf("temp dir %v is not in the expected location: %v", gotDir, os.TempDir())
			}

			// check proper creation of the dir
			_, err = os.Stat(gotDir)
			if err != nil {
				t.Errorf("unexpected error when inspecting the expected temp dir at %v: %v", gotDir, err)
			}

			// check proper cleanup of the dir
			if err := gotCleanup(); err != nil != tt.wantErr {
				t.Errorf("cleanup returned an error: %v, wantErr %v", err, tt.wantErr)
			}
			_, err = os.Stat(gotDir)
			if !errors.Is(err, fs.ErrNotExist) {
				t.Errorf("expected %v to have been deleted: %v", gotDir, err)
			}
		})
	}
}

func TestWithSubDirs(t *testing.T) {
	type args struct {
		subDirCount int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "create temp dir with 1 sub dirs",
			args:    args{subDirCount: 1},
			wantErr: false,
		},
		{
			name:    "create temp dir with 4 sub dirs",
			args:    args{subDirCount: 4},
			wantErr: false,
		},
		{
			name:    "create temp dir with 0 sub dirs",
			args:    args{subDirCount: 0},
			wantErr: true,
		},
		{
			name:    "create temp dir with negative sub dirs",
			args:    args{subDirCount: -1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithSubDirs(tt.args.subDirCount)

			dir, cleanup, err := TempDir(option)
			if err != nil {
				// check expected error
				if tt.wantErr {
					// pass test, no further processing possible
					return
				}
				t.Fatalf("TempDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer failOnError(t, cleanup)

			// check if the expected count of sub dirs has been created
			subDirCount := -1 // offset, since root dir is counted, too
			err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() {
					subDirCount += 1
				}
				return nil
			})
			if err != nil {
				t.Errorf("error encountered while walking through the temp dir tree: %v", err)
			}

			if tt.args.subDirCount != subDirCount {
				t.Errorf("expected %d sub dirs, got %d dirs", tt.args.subDirCount, subDirCount)
			}
		})
	}
}
