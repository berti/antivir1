package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const VirusMarkOffset int = 3
const VirusMarkEndOffset int = 5
const VirusGenerationOffset int = 5

func isComFile(path string, info os.FileInfo) bool {
	return !info.IsDir() && strings.EqualFold(filepath.Ext(path), ".com")
}

func isInfected(path string, info os.FileInfo) error {
	virusMark := []byte{0x49, 0x56}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	fileMark := content[VirusMarkOffset:VirusMarkEndOffset]
	if !bytes.Equal(virusMark, fileMark) {
		// Not infected
		return nil
	}

	fileGeneration := int(content[VirusGenerationOffset])
	fmt.Printf("Infected file found: %s, generation %d\n", path, fileGeneration)

	return nil
}

func processFile(path string, info os.FileInfo, err error) error {
	if isComFile(path, info) {
		return isInfected(path, info)
	}
	return nil
}

func findViruses(root string) {
	fmt.Printf("Finding viruses in %s\n", root)
	err := filepath.Walk(root, processFile)
	if err != nil {
		panic(err)
	}
}

func main() {
	var root string = "."

	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	findViruses(root)
}
