package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("get file info error: %w", err)
	}

	if !srcInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := srcInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	var bytesToCopy int64
	if limit == 0 {
		bytesToCopy = fileSize - offset
	} else {
		bytesToCopy = limit
		if offset+bytesToCopy > fileSize {
			bytesToCopy = fileSize - offset
		}
	}

	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek file error: %w", err)
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create destination file error: %w", err)
	}
	defer dstFile.Close()

	bar := progressbar.NewOptions64(
		bytesToCopy,
		progressbar.OptionSetDescription("Copying..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)

	reader := io.LimitReader(srcFile, bytesToCopy)
	multiReader := io.TeeReader(reader, bar)

	_, err = io.Copy(dstFile, multiReader)
	if err != nil {
		return fmt.Errorf("copy data error: %w", err)
	}

	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("sync destination file error: %w", err)
	}

	return nil
}
