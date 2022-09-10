package build

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func DirCopy(src_p string, des_p string) error {
	err := filepath.Walk(src_p, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		opath := filepath.Join(des_p, strings.TrimPrefix(path, src_p))

		if info.IsDir() {
			os.MkdirAll(opath, info.Mode().Perm())
			return nil
		}

		op, oerr := os.Open(path)
		if oerr != nil {
			return oerr
		}
		defer op.Close()

		dp, derr := os.Create(opath)
		if derr != nil {
			return derr
		}
		defer dp.Close()

		dp.Chmod(info.Mode())

		_, cerr := io.Copy(dp, op)
		return cerr
	})

	return err
}
