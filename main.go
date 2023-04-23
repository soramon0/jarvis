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
	fmt.Println("Projct name:", projectName)
	if err := createDirectory(projectName); err != nil {
		return err
	}
	return nil
}

func createDirectory(pathname string) error {
	file, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !file.IsDir() {
		return fmt.Errorf("%s is not a directory", file.Name())
	}
	ok, err := isDirEmpty(pathname)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%s is not an empty directory", pathname)
	}
	return nil
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases

}
