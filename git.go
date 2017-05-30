package main

import (
	"errors"
	"io/ioutil"
	"os/exec"

	"golang.org/x/crypto/ssh"

	"os"

	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func getAuth() (*gitssh.PublicKeys, error) {
	pem, err := ioutil.ReadFile(cfg.SSH.PrivateKey)
	if err != nil {
		return nil, errors.New("couldn't read private key '" + cfg.SSH.PrivateKey + "'")
	}

	signer, err := ssh.ParsePrivateKey(pem)
	if err != nil {
		return nil, err
	}

	auth := &gitssh.PublicKeys{
		User:   cfg.SSH.Username,
		Signer: signer,
	}

	return auth, nil
}

// gitClone uses go-git functions to clone.
func gitClone(url, path string) error {
	auth, err := getAuth()
	if err != nil {
		return err
	}

	o := git.CloneOptions{
		URL:  url,
		Auth: auth,
	}

	path = filepath.Join(cfg.Repos, path)
	_, err = git.PlainClone(path, false, &o)
	if err != nil {
		return err
	}

	return nil
}

// gitPull updates a repository via go-git functions.
func gitPull(path string) error {
	auth, err := getAuth()
	if err != nil {
		return err
	}

	o := git.PullOptions{
		Auth: auth,
	}

	path = filepath.Join(cfg.Repos, path)
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	err = repo.Pull(&o)
	return err
}

// gitExternalClone clones via whichever command line git is installed in the user's path.
// Use with anything which go-git can't handle, like VSTS.
func gitExternalClone(url, path string) error {
	cmd := exec.Command("git", "-C", cfg.Repos, "clone", url, path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// gitExternalPull updates via the installed git command.
func gitExternalPull(path string) error {
	path = filepath.Join(cfg.Repos, path)
	cmd := exec.Command("git", "-C", path, "pull")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
