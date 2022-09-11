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
			_ = op.Close()
		}(op)

		dp, derr := os.Create(opath)
		if derr != nil {
			return derr
		}
		defer func(dp *os.File) {
			_ = dp.Close()
		}(dp)

		_ = dp.Chmod(info.Mode())

		_, cerr := io.Copy(dp, op)
		return cerr
	})

	return err
}
