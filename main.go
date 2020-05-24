package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const VirusMarkOffset int = 3
const VirusMarkEndOffset int = 5
const VirusGenerationOffset int = 5
const OriginalCodeOffset int = 0x34f

const VirusNotFound int = -1

func isComFile(path string, info os.FileInfo) bool {
	return !info.IsDir() && strings.EqualFold(filepath.Ext(path), ".com")
}

func isInfected(path string, info os.FileInfo) (int, error) {
	virusMark := []byte{0x49, 0x56}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return VirusNotFound, err
	}

	fileMark := content[VirusMarkOffset:VirusMarkEndOffset]
	if !bytes.Equal(virusMark, fileMark) {
		// Not infected
		return VirusNotFound, nil
	}

	fileGeneration := int(content[VirusGenerationOffset])
	fmt.Printf("Infected file found: %s, generation %d\n", path, fileGeneration)

	return fileGeneration, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func removeVirus(path string, info os.FileInfo) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	offset := min(OriginalCodeOffset, int(info.Size()))
	originalCode := content[offset:]
	ioutil.WriteFile(path, originalCode, info.Mode())

	fmt.Println("Virus removed")

	return nil
}

func createFileProcessor(remove bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if isComFile(path, info) {
			gen, err := isInfected(path, info)
			if err != nil {
				return err
			}
			if gen != VirusNotFound && remove {
				return removeVirus(path, info)
			}
		}
		return nil
	}
}

func findViruses(root string, remove bool) {
	fmt.Printf("Finding viruses in %s\n", root)
	if remove {
		fmt.Println("Found viruses will be removed")
	}

	processFile := createFileProcessor(remove)
	err := filepath.Walk(root, processFile)
	if err != nil {
		panic(err)
	}
}

func main() {
	pathPtr := flag.String("p", ".", "path to scan")
	removePtr := flag.Bool("r", false, "remove virus payload from infected files")

	flag.Parse()

	findViruses(*pathPtr, *removePtr)
}
