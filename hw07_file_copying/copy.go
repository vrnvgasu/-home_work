package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	inFile, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open inFile: %w", err)
	}
	defer AppClose(inFile)

	outFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create outFile: %w", err)
	}
	defer AppClose(outFile)

	stat, err := inFile.Stat()
	if err != nil {
		return fmt.Errorf("stat inFile: %w", err)
	}

	if err = checkOffset(offset, stat.Size()); err != nil {
		return fmt.Errorf("checkOffset: %w", err)
	}
	if _, err = inFile.Seek(offset, 0); err != nil {
		return fmt.Errorf("seek inFile: %w", err)
	}

	if err = copyFile(inFile, outFile, prepareLimit(limit, stat.Size())); err != nil {
		return fmt.Errorf("copyFile: %w", err)
	}

	return nil
}

func copyFile(in io.Reader, out io.Writer, limit int64) error {
	bar := pb.New(int(limit))
	bar.Start()

	reader := bar.NewProxyReader(in)

	if _, err := io.CopyN(out, reader, limit); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("copyN: %w", err)
	}

	bar.Finish()

	return nil
}

func prepareLimit(limit int64, fileSize int64) (newLimit int64) {
	if limit == 0 || limit > fileSize {
		return fileSize
	}

	return limit
}

func checkOffset(offset int64, fileSize int64) error {
	switch {
	case fileSize < 0:
		return fmt.Errorf("offset is negative: %d", offset)
	case offset > fileSize:
		return ErrOffsetExceedsFileSize
	default:
		return nil
	}
}
