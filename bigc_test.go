package bigc_test

import (
	"fmt"

	"github.com/mathlava/bigc"
)

func Example() {
	ex, _ := bigc.ParseString("(1+2i)/(3-4i)")
	fmt.Println(ex.String())
}
