package main

import (
	"flag"
	"fmt"
)

var (
	// version of nomad-var-dirsync being run
	version = "dev"
)

func main() {
	showVersion := flag.Bool("version", false, "Display the version of nomad-var-dirsync and exit")
	flag.Parse()

	// Print version if flag is provided
	if *showVersion {
		fmt.Println("nomad-var-dirsync version:", version)

		return
	}
}
