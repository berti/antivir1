package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

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

	offset := min(originalCodeOffset, int(info.Size()))
	originalCode := content[offset:]
	ioutil.WriteFile(path, originalCode, info.Mode())

	fmt.Printf("Virus removed: %s\n", path)

	return nil
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
