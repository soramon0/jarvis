package project

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Project struct {
	name string
	// friendlyName returns either directory name or
	// name of current directory if d.name is ".".
	// Used mainly while printing errors
	friendlyName string
	absPath      string
	pwd          string
}

func initProject(file *os.File) (*Project, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	name := file.Name()
	friendlyName := name
	absPath, err := filepath.Abs(friendlyName)
	if err != nil {
		return nil, err
	}
	friendlyName = filepath.Base(absPath)
	return &Project{name: name, friendlyName: friendlyName, absPath: absPath, pwd: pwd}, nil
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) FriendlyName() string {
	return p.friendlyName
}

func (p *Project) Pwd() string {
	return p.pwd
}

func (p *Project) AbsPath() string {
	return p.absPath
}

func (p *Project) ChdirToProjectDir() error {
	return os.Chdir(p.absPath)
}

func (p *Project) ChdirToPWD() error {
	return os.Chdir(p.pwd)
}

// newProject creates a directory to work in if it doesn't already exist.
// If directory exists but it's not empty an error is returend
func newProject(pathname string) (*Project, error) {
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
			return initProject(file)
		}
	}
	defer file.Close()

	p, err := initProject(file)
	if err != nil {
		return nil, err
	}
	if err := isEmpty(file, p.friendlyName); err != nil {
		return nil, err
	}

	return p, nil
}

// Clean removes directory and any children it cantains
func (p *Project) Clean() error {
	return os.RemoveAll(p.name)
}

// isEmpty returns an error if directory is not a directory
// or if it is not an empty one
func isEmpty(file *os.File, friendlyName string) error {
	f, err := file.Stat()
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is not a directory", friendlyName)
	}

	fileNames, err := file.Readdirnames(1)
	if err != nil && err != io.EOF {
		return err
	}
	if len(fileNames) > 0 {
		return fmt.Errorf("%s is not an empty directory", friendlyName)
	}

	return nil
}
