package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt(mustfilled bool) (string, error) {
	var result string
	var ioerr error
	red := bufio.NewReader(os.Stdin)

	for {
		result, ioerr = red.ReadString('\n')

		if ioerr != nil {
			return "", ioerr
		}

		if result == "\n" && mustfilled {
			fmt.Println("required field")
			continue
		}

		result = strings.ReplaceAll(result, "\n", "")
		break
	}

	return result, nil
}
