package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt(mustfilled bool) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var result string

	for {
		result, ioerr := reader.ReadString('\n')

		if ioerr != nil {
			return "", ioerr
		}

		// Trim the newline character from the input
		result = strings.TrimSuffix(result, "\n")

		if result != "" || !mustfilled {
			break
		}

		fmt.Println("required field")
	}

	return result, nil
}
