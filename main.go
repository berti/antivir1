package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isInfectable(path string, info os.FileInfo) bool {
	return !info.IsDir() && strings.EqualFold(filepath.Ext(path), ".com")
}

func isInfected(path string, info os.FileInfo) error {
	src := []byte{0xEB, 0x14, 0x90, 0x49, 0x56, 0x01, 0x2A, 0x2E, 0x43, 0x4F, 0x4D, 0x00}

	size := info.Size()
	if size < int64(len(src)) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dst := make([]byte, len(src))
	_, err = file.Read(dst)
	if err != nil {
		return err
	}

	if bytes.Equal(src, dst) {
		fmt.Printf("%s is infected!\n", path)
	}

	return nil
}

func processFile(path string, info os.FileInfo, err error) error {
	if isInfectable(path, info) {
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
