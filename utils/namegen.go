package utils

import (
	"encoding/base64"
	"strconv"
	"strings"
)

var FilterStr []string = []string{
	"/", `\`, ">", "<", ".", ":", `"`, "|", "?", "*",
}

func MakeUniqueFileName(title string, stampdata int) string {
	var fileTitle string

	b64name := base64.RawStdEncoding.EncodeToString([]byte(title))
	stampsd := strconv.Itoa(stampdata)

	for _, ct := range FilterStr {
		b64name = strings.ReplaceAll(b64name, ct, "_")
	}

	if len(b64name) > 6 {
		b64name = b64name[:6]
	}

	if len(stampsd) > 4 {
		stampsd = stampsd[:4]
	}

	fileTitle = stampsd + b64name

	return fileTitle
}
