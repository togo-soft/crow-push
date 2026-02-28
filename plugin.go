package main

import (
	"errors"
	"fmt"

	"codefloe.com/actions/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// Run performs the push operation for all enabled platforms
func Run(cfg *Config) error {
	repo, err := git.PlainOpen(cfg.Workspace)
	if err != nil {
		return fmt.Errorf("open repository: %w", err)
	}

	var errors []error
	for _, platform := range cfg.Platforms {
		if !platform.Enabled {
			continue
		}

		common.Info("pushing to platform: %s (%s)", platform.Name, platform.URL)
		if err := pushToPlatform(repo, platform); err != nil {
			common.Error("failed to push to %s: %v", platform.Name, err)
			errors = append(errors, fmt.Errorf("%s: %w", platform.Name, err))
		} else {
			common.Info("successfully pushed to platform: %s", platform.Name)
		}
	}

	if len(errors) > 0 {
		errMsg := "push to multiple platforms failed:\n"
		for _, e := range errors {
			errMsg += "  - " + e.Error() + "\n"
		}
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

// pushToPlatform pushes all refs to a specific platform
func pushToPlatform(repo *git.Repository, platform *Platform) error {
	remoteName := remoteName(platform)

	// Remove existing remote if present
	if err := repo.DeleteRemote(remoteName); err != nil && !errors.Is(err, git.ErrRemoteNotFound) {
		return fmt.Errorf("delete remote: %w", err)
	}

	// Create remote
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{platform.URL},
	})
	if err != nil {
		return fmt.Errorf("create remote: %w", err)
	}

	// Build authentication
	auth, err := buildAuth(platform)
	if err != nil {
		return fmt.Errorf("build auth: %w", err)
	}

	// Push all refs
	refSpecs := []config.RefSpec{
		"refs/heads/*:refs/heads/*",
		"refs/tags/*:refs/tags/*",
	}

	err = repo.Push(&git.PushOptions{
		RemoteName:      remoteName,
		RefSpecs:        refSpecs,
		Auth:            auth,
		Force:           true,
		InsecureSkipTLS: true,
	})

	// NoErrAlreadyUpToDate is success
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}

	return err
}

// buildAuth returns the appropriate authentication method for a platform
func buildAuth(platform *Platform) (transport.AuthMethod, error) {
	if platform.SSHKey != "" {
		publicKeys, err := ssh.NewPublicKeys("git", []byte(platform.SSHKey), platform.SSHKeyPassphrase)
		if err != nil {
			return nil, fmt.Errorf("parse ssh key: %w", err)
		}

		// Set HostKeyCallback to InsecureIgnoreHostKey to disable host key verification in CI environments
		publicKeys.HostKeyCallback = gossh.InsecureIgnoreHostKey()

		return publicKeys, nil
	}

	if platform.Username != "" && platform.Token != "" {
		return &http.BasicAuth{
			Username: platform.Username,
			Password: platform.Token,
		}, nil
	}

	if platform.Username != "" && platform.Password != "" {
		return &http.BasicAuth{
			Username: platform.Username,
			Password: platform.Password,
		}, nil
	}

	return nil, fmt.Errorf("no auth method available")
}

// remoteName returns the remote name for a platform
func remoteName(platform *Platform) string {
	if platform.RemoteName != "" {
		return platform.RemoteName
	}
	return "push-plugin-" + platform.Name
}
