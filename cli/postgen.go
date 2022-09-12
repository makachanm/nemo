package cli

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GeneratePost() {
	wd, _ := os.Getwd()

	year := strconv.Itoa(time.Now().Year())
	month := strconv.Itoa(int(time.Now().Month()))
	day := strconv.Itoa(time.Now().Day())
	hour := strconv.Itoa(time.Now().Hour())
	min := strconv.Itoa(time.Now().Minute())

	timest := fmt.Sprintf("year=%v,month=%v,day=%v,hour=%v,min=%v", year, month, day, hour, min)
	filename := fmt.Sprintf("%v%v%v%v%v", year, month, day, hour, min)

	postctx := "$[title Title]\n$[summary Summary of Post]\n$[timestamp(" + timest + ")]\n==========\n"

	dir := wd + "/post/" + filename + ".ps"
	os.WriteFile(dir, []byte(postctx), 0777)
}
