package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GeneratePost() {
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(now.Day())
	hour := strconv.Itoa(now.Hour())
	min := strconv.Itoa(now.Minute())

	timest := fmt.Sprintf("year=%v,month=%v,day=%v,hour=%v,min=%v", year, month, day, hour, min)

	postFileName := fmt.Sprintf("%v%v%v%v%v", year, month, day, hour, min)

	postctx := fmt.Sprintf(`
$[title Title]
$[summary Summary of Post]
$[timestamp(%v)]
$[tag Tags]
==========`, timest)

	postFilePath := filepath.Join("post", postFileName+".ps")

	postFile, err := os.Create(postFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(postFile *os.File) {
		err := postFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(postFile)

	if _, err := io.Copy(postFile, bytes.NewBuffer([]byte(postctx))); err != nil {
		log.Fatal(err)
	}
}
