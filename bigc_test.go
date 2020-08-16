package bigc_test

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mathlava/bigc"
)

func Example() {
	ex, _ := bigc.ParseString("(1+2i)/(3-4i)")
	fmt.Println(ex.String())
}
