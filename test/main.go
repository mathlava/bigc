package main

import (
	"fmt"

	"github.com/mathlava/bigc"
)

func main() {
	for {
		var in string
		fmt.Scanln(&in)
		ex, err := bigc.ParseString("12 + 31i")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ex.String())
	}
}
