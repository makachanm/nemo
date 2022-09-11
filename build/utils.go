package build

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func DirCopy(srcP string, desP string) error {
	err := filepath.Walk(srcP, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		opath := filepath.Join(desP, strings.TrimPrefix(path, srcP))

		if info.IsDir() {
			_ = os.MkdirAll(opath, info.Mode().Perm())
			return nil
		}

		op, oerr := os.Open(path)
		if oerr != nil {
			return oerr
		}
		defer func(op *os.File) {
			err := op.Close()
			if err != nil {

			}
		}(op)

		dp, derr := os.Create(opath)
		if derr != nil {
			return derr
		}
		defer func(dp *os.File) {
			err := dp.Close()
			if err != nil {

			}
		}(dp)

		_ = dp.Chmod(info.Mode())

		_, cerr := io.Copy(dp, op)
		return cerr
	})

	return err
}
