package util

import (
	"fmt"
	"os"
)

// PrintlnError prints error to the stderr; if failed — to stdout.
func PrintlnError(err error) {
	if err == nil {
		return
	}

	if _, e := fmt.Fprintln(os.Stderr, err.Error()); e != nil {
		fmt.Println(err.Error())
	}
}
