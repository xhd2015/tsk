package main

import (
	"fmt"
	"os"

	"github.com/xhd2015/tsk/tskcli"
)

func main() {
	if err := tskcli.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}