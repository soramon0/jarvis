package project

import (
	"fmt"
	"os"
	"os/exec"
)

func setupGitRepo(projectName string) error {
	cmd := exec.Command("git", "init")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}

	readmeBytes := []byte(fmt.Sprintf("# %s\n", projectName))
	if err := os.WriteFile("README.md", readmeBytes, 0644); err != nil {
		return err
	}
	cmd = exec.Command("git", "add", "README.md")
	if err := cmd.Run(); err != nil {
		return err
	}

	return exec.Command("git", "commit", "-m", "init project").Run()
}

func pushToGithub(projectName string, isRepoPrivate bool) error {
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("gh command is not found. please install github cli to push repository to github")
	}
	cmd := exec.Command(ghPath, "repo", "create", projectName, "--source=.", "--remote=upstream", "--push")
	cmd.Stdout = os.Stdout

	if isRepoPrivate {
		cmd.Args = append(cmd.Args, "--private")
	} else {
		cmd.Args = append(cmd.Args, "--public")
	}
	return cmd.Run()
}

func removeRemoteRepository(projectName string) error {
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("gh command is not found. please install github cli to delete repository from github")
	}
	// TODO(sora): check repo exists
	// TODO(sora): check permissions first
	cmd := exec.Command(ghPath, "repo", "delete", projectName, "--yes")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to delete remote repo %v", err)
	}
	return nil
}
