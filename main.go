package main

import (
	"flag"
	"fmt"

	"github.com/ceaser/pac/internal/version"
)

func main() {
	flag.Parse()
	version.ShowVersion()
	fmt.Println("vim-go")
}
