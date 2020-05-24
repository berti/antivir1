package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func isComFile(path string, info os.FileInfo) bool {
	return !info.IsDir() && strings.EqualFold(filepath.Ext(path), ".com")
}

func isInfected(path string, info os.FileInfo) (int, error) {
	virusMark := []byte{0x49, 0x56}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return virusNotFound, err
	}

	fileMark := content[virusMarkOffset:virusMarkEndOffset]
	if !bytes.Equal(virusMark, fileMark) {
		// Not infected
		return virusNotFound, nil
	}

	fileGeneration := int(content[virusGenerationOffset])
	fmt.Printf("Infected file found: %s, generation %d\n", path, fileGeneration)

	return fileGeneration, nil
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
			if gen != virusNotFound {
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
