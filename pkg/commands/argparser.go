package commands

import (
	"os"

	"github.com/soramon0/jarvis/pkg/project"
)

func (c *Cmd) ParseCreateProjectArgs() (*project.CreateArgs, error) {
	shouldPushToGithub := c.Bool("gpush", false, "Push local repository to the Github using Github CLI")
	isRepoPrivate := c.Bool("gprivate", false, "Make the new repository private")
	projectName := c.String("p", "", "Project name")
	if err := c.Parse(os.Args[2:]); err != nil {
		return nil, err
	}

	return &project.CreateArgs{
		ProjectName:        *projectName,
		ShouldPushToGithub: *shouldPushToGithub,
		IsRepoPrivate:      *isRepoPrivate,
	}, nil
}

func (c *Cmd) ParseDeleteProjectArgs() (*project.DeleteArgs, error) {
	remoteRemove := c.Bool("remote", false, "delete remote repository from github")
	remoteOnly := c.Bool("remoteOnly", false, "delete only github remote repository")
	projectName := c.String("p", "", "Project name")
	if err := c.Parse(os.Args[2:]); err != nil {
		return nil, err
	}

	return &project.DeleteArgs{
		ProjectName:  *projectName,
		RemoteRemove: *remoteRemove,
		RemoteOlny:   *remoteOnly,
	}, nil
}
