package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	tmpDir := os.TempDir()

	tests := []struct {
		name        string
		from        string
		to          string
		offset      int64
		limit       int64
		wantErr     bool
		errType     error
		compareWith string
	}{
		{
			name:        "offset0_limit0",
			from:        "testdata/input.txt",
			to:          filepath.Join(tmpDir, "test1.txt"),
			offset:      0,
			limit:       0,
			wantErr:     false,
			compareWith: "testdata/out_offset0_limit0.txt",
		},
		{
			name:        "offset0_limit10",
			from:        "testdata/input.txt",
			to:          filepath.Join(tmpDir, "test2.txt"),
			offset:      0,
			limit:       10,
			wantErr:     false,
			compareWith: "testdata/out_offset0_limit10.txt",
		},
		{
			name:        "_offset0_limit1000",
			from:        "testdata/input.txt",
			to:          filepath.Join(tmpDir, "test3.txt"),
			offset:      0,
			limit:       1000,
			wantErr:     false,
			compareWith: "testdata/out_offset0_limit1000.txt",
		},
		{
			name:        "out_offset0_limit10000",
			from:        "testdata/input.txt",
			to:          filepath.Join(tmpDir, "test4.txt"),
			offset:      0,
			limit:       10000,
			wantErr:     false,
			compareWith: "testdata/out_offset0_limit10000.txt",
		},
		{
			name:        "out_offset100_limit1000",
			from:        "testdata/input.txt",
			to:          filepath.Join(tmpDir, "test5.txt"),
			offset:      100,
			limit:       1000,
			wantErr:     false,
			compareWith: "testdata/out_offset100_limit1000.txt",
		},
		{
			name:        "offset6000_limit1000",
			from:        "testdata/input.txt",
			to:          filepath.Join(tmpDir, "test6.txt"),
			offset:      6000,
			limit:       1000,
			wantErr:     false,
			compareWith: "testdata/out_offset6000_limit1000.txt",
		},
		{
			name:    "offset exceeds file size",
			from:    "testdata/input.txt",
			to:      filepath.Join(tmpDir, "test7.txt"),
			offset:  1000000,
			limit:   1000,
			wantErr: true,
			errType: ErrOffsetExceedsFileSize,
		},
		{
			name:    "non-existent source file",
			from:    "testdata/nonexistent.txt",
			to:      filepath.Join(tmpDir, "test8.txt"),
			offset:  0,
			limit:   1000,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if _, err := os.Stat(tt.to); err == nil {
					os.Remove(tt.to)
				}
			}()

			err := Copy(tt.from, tt.to, tt.offset, tt.limit)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Copy error: nil")
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("Copy error:  %v, want %v", err, tt.errType)
				}
			} else {
				if err != nil {
					t.Errorf("Copy unexpected error: %v", err)
					return
				}

				if _, err := os.Stat(tt.to); os.IsNotExist(err) {
					t.Errorf("Copy error: destination file not created")
					return
				}

				if tt.compareWith != "" {
					expected, err := os.ReadFile(tt.compareWith)
					if err != nil {
						t.Errorf("Read expected file error: %v", err)
						return
					}

					actual, err := os.ReadFile(tt.to)
					if err != nil {
						t.Errorf("Read actual file error: %v", err)
						return
					}

					if string(actual) != string(expected) {
						//fmt.Println("--------------------------------------------")
						//fmt.Println(actual)
						//fmt.Println("--------------------------------------------")
						//fmt.Println(expected)
						//fmt.Println("--------------------------------------------")
						t.Errorf("Copy content mismatch")
					}
				}
			}
		})
	}
}
