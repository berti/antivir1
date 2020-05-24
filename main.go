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

	fmt.Printf("Virus removed: %s\n", path)

	return nil
}

func find(root string) []string {
	fmt.Printf("Finding viruses in %s\n", root)

	var infectedFiles []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if isComFile(path, info) {
			gen, err := isInfected(path, info)
			if err != nil {
				return err
			}
			if gen != VirusNotFound {
				infectedFiles = append(infectedFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return infectedFiles
}

func remove(files []string) {
	fmt.Println("Removing found viruses")

	for _, file := range files {
		fileInfo, err := os.Lstat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while getting file info: %s", file)
			continue
		}

		err = removeVirus(file, fileInfo)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while removing virus: %s", file)
			continue
		}
	}
}

func main() {
	pathPtr := flag.String("p", ".", "path to scan")
	removePtr := flag.Bool("r", false, "remove virus payload from infected files")

	flag.Parse()

	infectedFiles := find(*pathPtr)
	if *removePtr {
		remove(infectedFiles)
	}
}
