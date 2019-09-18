package main

import (
	"fmt"
	"log"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
)

type GitHubDownloader struct {
	path string
}

func GitHub(path string) *GitHubDownloader {
	return &GitHubDownloader{path: path}
}

func (gh *GitHubDownloader) Update(folder, name string) error {
	log.Printf("Updating github repo %s", name)

	url := fmt.Sprintf("https://github.com/%s.git", name)
	path := filepath.Join(gh.path, folder)

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: url,
	})

	if err == git.ErrRepositoryAlreadyExists {
		return nil
	}

	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		Force: true,
	})

	if err == git.NoErrAlreadyUpToDate {
		return nil
	}

	return err
}
