package build

import (
	"fmt"
	"io"
	"os"
)

func DirCopy(src string, dst string) error {
	// Get properties of source dir
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination dir
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(src)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		srcPointer := src + "/" + obj.Name()
		dstPointer := dst + "/" + obj.Name()

		if obj.IsDir() {
			// Create sub-directories - recursively
			err = DirCopy(srcPointer, dstPointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Perform copy
			err = FileCopy(srcPointer, dstPointer)
			if err != nil {
				fmt.Println(err)
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
			panic(err)
		}
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
			panic(err)
		}
	}(destination)
	_, err = io.Copy(destination, source)
	return err
}
