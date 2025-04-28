package main

import (
	"flag"
	"fmt"
	"glstatus/components"
	"os"
	"time"
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
	for {
		var status string
		for _, mod := range modules {
			output := mod.Producer(mod.Argument)
			if output == "" {
				output = components.UnknownStr
			}
			status += fmt.Sprintf(mod.Format, output)
		}

		if !silentMode {
			fmt.Print("\r" + status)
		} else {
			fmt.Println(status)
		}

		time.Sleep(components.UpdateInterval)
	}
}
