package main

import (
	"flag"
	"fmt"
	"os"
)

var version bool
var silentMode bool

func init() {
	flag.BoolVar(&version, "v", false, "Print the package version")
	flag.BoolVar(&silentMode, "s", false, "Silent mode")
}

func usage() {
	fmt.Printf("usage: %s [flags]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if version {
		fmt.Printf("glstatus-%s\n", "0.0.1")
		return
	}

	if !silentMode {
		fmt.Println("Print hardware information")
	} else {
		fmt.Println("X11 EWMH root window property not implemented yet")
	}
}
