package main

import (
	"flag"
)

func main() {
	pathPtr := flag.String("p", ".", "path to scan")
	removePtr := flag.Bool("r", false, "remove virus payload from infected files")

	flag.Parse()

	infectedFiles := find(*pathPtr)
	if *removePtr {
		remove(infectedFiles)
	}
}
