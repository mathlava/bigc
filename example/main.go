package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mathlava/bigc"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ex, err := bigc.ParseString(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ex.String())
	}
}
