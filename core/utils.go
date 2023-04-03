package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const bufferSize = 32 * 1024

func DirCopy(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source directory properties: %w", err)
	}
	if src == dst {
		return fmt.Errorf("source and destination directories are the same")
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	directory, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source directory: %w", err)
	}
	defer func(directory *os.File) {
		_ = directory.Close()
	}(directory)
	objects, err := directory.Readdir(-1)
	if err != nil {
		return fmt.Errorf("failed to read source directory contents: %w", err)
	}
	for _, obj := range objects {
		srcPointer := filepath.Join(src, obj.Name())
		dstPointer := filepath.Join(dst, obj.Name())
		if obj.IsDir() {
			if err := DirCopy(srcPointer, dstPointer); err != nil {
				return fmt.Errorf("failed to copy sub-directory %s: %w", srcPointer, err)
			}
		} else {
			if err := FileCopy(srcPointer, dstPointer); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", srcPointer, err)
			}
		}
	}
	return nil
}

func FileCopy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file properties: %w", err)
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)
	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)
	buf := make([]byte, bufferSize)
	_, err = io.CopyBuffer(destination, source, buf)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}
	return nil
}
