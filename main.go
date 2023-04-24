package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("project name is required")
	}
	projectName := os.Args[len(os.Args)-1]
	dir, err := createDirectory(projectName)
	if err != nil {
		return err
	}
	defer dir.Close()

	cmd := exec.Command("git", "init", dir.Name())
	if err := cmd.Run(); err != nil {
		if err := dir.Clean(); err != nil {
			fmt.Printf("failed to clean %s directory\n", dir.FriendlyName())
		}
		return err
	}

	return nil
}

type directory struct {
	*os.File
	friendlyName string
}

func newDirectory(file *os.File) *directory {
	return &directory{File: file}
}

// createDirectory creates a directory to work in if it doesn't already exist.
// If directory exists but it's not empty an error is returend
func createDirectory(pathname string) (*directory, error) {
	file, err := os.Open(filepath.Clean(pathname))
	if err != nil {
		// only move forward if file does not exist
		if !os.IsNotExist(err) {
			return nil, err
		}
		// Stop execution whether we create the directory or fail
		if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
			return nil, err
		} else {
			file, err := os.Open(pathname)
			if err != nil {
				return nil, err
			}
			return newDirectory(file), nil
		}
	}

	dir := newDirectory(file)
	if err := dir.IsEmpty(); err != nil {
		return nil, err
	}

	return dir, nil
}

// Clean removes directory and any children it cantains
func (d *directory) Clean() error {
	if d.File == nil {
		return os.ErrInvalid
	}
	return os.RemoveAll(d.File.Name())
}

// FriendlyName returns either directory name or
// name of current directory if d.name is ".".
// Used mainly while printing errors
func (d *directory) FriendlyName() string {
	if d.friendlyName != "" {
		return d.friendlyName
	}
	d.friendlyName = d.File.Name()
	if d.friendlyName == "." {
		absPath, err := filepath.Abs(d.friendlyName)
		if err == nil {
			d.friendlyName = filepath.Base(absPath)
		}
	}
	return d.friendlyName
}

// IsEmpty returns an error if directory is not a directory
// or if it is not an empty one
func (d *directory) IsEmpty() error {
	if d.File == nil {
		return os.ErrInvalid
	}
	f, err := d.File.Stat()
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is not a directory", d.FriendlyName())
	}

	fileNames, err := d.File.Readdirnames(1)
	if err != nil && err != io.EOF {
		return err
	}
	if len(fileNames) > 0 {
		return fmt.Errorf("%s is not an empty directory", d.FriendlyName())
	}

	return nil
}
