package utils

import (
	"encoding/base64"
	"strings"
)

var FilterStr []string = []string{
	"/", `\`, ">", "<", ".", ":", `"`, "|", "?", "*",
}

func MakeUniqueFileName(title string) string {
	var fileTitle string

	fname_k := base64.RawStdEncoding.EncodeToString([]byte(title))

	for _, ct := range FilterStr {
		fname_k = strings.ReplaceAll(fname_k, ct, "_")
	}

	if len(fname_k) > 16 {
		fileTitle = fname_k[:16]
	} else {
		fileTitle = fname_k
	}

	return fileTitle
}
