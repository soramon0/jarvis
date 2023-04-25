package project

import (
	"fmt"
	"os"
	"path/filepath"
)

type CreateArgs struct {
	ProjectName        string
	ShouldPushToGithub bool
	IsRepoPrivate      bool
}

func Create(args *CreateArgs) (*Project, error) {
	if args.ProjectName == "" {
		return nil, fmt.Errorf("project name is required")
	}

	p, err := newProject(args.ProjectName)
	if err != nil {
		return nil, err
	}

	if err := p.ChdirToProjectDir(); err != nil {
		return nil, err
	}
	defer p.ChdirToPWD()

	if err := setupGitRepo(args.ProjectName); err != nil {
		if err := p.Clean(); err != nil {
			fmt.Printf("failed to clean %s directory. %v \n", p.friendlyName, err)
		}
		return nil, err
	}

	if args.ShouldPushToGithub {
		if err := pushToGithub(args.ProjectName, args.IsRepoPrivate); err != nil {
			if err := p.Clean(); err != nil {
				fmt.Printf("failed to clean %s directory. %v \n", p.friendlyName, err)
			}
			return nil, err
		}
	}

	return p, err
}

type DeleteArgs struct {
	ProjectName  string
	RemoteRemove bool
	RemoteOlny   bool
}

func Delete(args *DeleteArgs) (*Project, error) {
	if args.ProjectName == "" {
		return nil, fmt.Errorf("project name is required")
	}

	var p *Project

	if !args.RemoteOlny {
		file, err := os.Open(filepath.Clean(args.ProjectName))
		if err != nil {
			return nil, err
		}

		f, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if !f.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", f.Name())
		}

		project, err := initProject(file)
		if err != nil {
			return nil, err
		}
		p = project

		if err := p.Clean(); err != nil {
			return nil, err
		}
	}

	if args.RemoteRemove || args.RemoteOlny {
		if err := removeRemoteRepository(args.ProjectName); err != nil {
			return nil, err
		}
	}

	return p, nil
}
