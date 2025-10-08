package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type testDto struct {
	name        string
	from        string
	to          string
	offset      int64
	limit       int64
	wantErr     bool
	errType     error
	compareWith string
}

func TestCopyOk(t *testing.T) {
	tmpDir := os.TempDir()

	tests := []testDto{
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
	}

	processTest(t, tests)
}

func TestCopyError(t *testing.T) {
	tmpDir := os.TempDir()

	tests := []testDto{
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

	processTest(t, tests)
}

func processTest(t *testing.T, tests []testDto) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupTestFile(t, tt.to)
			err := Copy(tt.from, tt.to, tt.offset, tt.limit)

			if tt.wantErr {
				assertError(t, err, tt.errType)
			} else {
				assertSuccess(t, err, tt.to, tt.compareWith)
			}
		})
	}
}

func cleanupTestFile(t *testing.T, filePath string) {
	t.Helper()
	t.Cleanup(func() {
		if _, err := os.Stat(filePath); err == nil {
			if removeErr := os.Remove(filePath); removeErr != nil {
				t.Logf("Failed to cleanup test file %s: %v", filePath, removeErr)
			}
		}
	})
}

func assertError(t *testing.T, err error, expectedErrType error) {
	t.Helper()
	if err == nil {
		t.Error("Copy error: nil")
		return
	}

	if expectedErrType != nil && !errors.Is(err, expectedErrType) {
		t.Errorf("Copy error: got %v, want %v", err, expectedErrType)
	}
}

func assertSuccess(t *testing.T, err error, destFile string, compareFile string) {
	t.Helper()
	if err != nil {
		t.Errorf("Copy unexpected error: %v", err)
		return
	}

	assertFileExists(t, destFile)

	if compareFile != "" {
		assertFilesEqual(t, destFile, compareFile)
	}
}

func assertFileExists(t *testing.T, filePath string) {
	t.Helper()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Copy error: destination file %s not created", filePath)
	}
}

func assertFilesEqual(t *testing.T, actualFile string, expectedFile string) {
	t.Helper()
	expected, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Errorf("Read expected file error: %v", err)
		return
	}

	actual, err := os.ReadFile(actualFile)
	if err != nil {
		t.Errorf("Read actual file error: %v", err)
		return
	}

	if string(actual) != string(expected) {
		t.Errorf("Copy content mismatch")
	}
}
