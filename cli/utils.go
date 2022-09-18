package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt(mustfilled bool) string {
	red := bufio.NewReader(os.Stdin)
	result, err := red.ReadString('\n')

	if err != nil {
		panic(err)
	}

	if result == "\n" && mustfilled {
		fmt.Println("required feild")
		return Prompt(mustfilled)
	}

	result = strings.ReplaceAll(result, "\n", "")

	return result
}
