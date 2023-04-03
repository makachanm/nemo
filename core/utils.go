package build

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func DirCopy(src string, dst string) error {
	// Get properties of source dir
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if src and dst directories are the same
	if src == dst {
		return fmt.Errorf("source and destination directories are the same")
	}

	// Create destination dir
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	directory, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(directory *os.File) {
		err := directory.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(directory)

	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		srcPointer := src + string(filepath.Separator) + obj.Name()
		dstPointer := dst + string(filepath.Separator) + obj.Name()

		if obj.IsDir() {
			// Create sub-directories - recursively
			err = DirCopy(srcPointer, dstPointer)
			if err != nil {
				return err
			}
		} else {
			// Perform copy
			err = FileCopy(srcPointer, dstPointer)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func FileCopy(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(destination)

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(destination, source, buf)
	return err
}
