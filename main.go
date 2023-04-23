package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	// create folder if empty
	if len(os.Args) < 2 {
		return fmt.Errorf("project name is required")
	}
	projectName := os.Args[len(os.Args)-1]
	if err := createDirectory(projectName); err != nil {
		return err
	}
	return nil
}

func createDirectory(pathname string) error {
	file, err := os.Open(pathname)
	if err != nil {
		// only move forward if file does not exist
		if !os.IsNotExist(err) {
			return err
		}
		// Stop execution whether we create the directory or fail
		if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
			return err
		} else {
			return nil
		}
	}
	defer file.Close()

	return isDirEmpty(file)
}

func isDirEmpty(file *os.File) error {
	f, err := file.Stat()
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is not a directory", file.Name())
	}

	fileNames, err := file.Readdirnames(1)
	if err != nil && err != io.EOF {
		return err
	}
	if len(fileNames) > 0 {
		return fmt.Errorf("%s is not an empty directory", file.Name())
	}

	return nil
}
