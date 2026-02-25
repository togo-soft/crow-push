package main

import (
	"encoding/json"
	"fmt"

	"codefloe.com/actions/common"
)

type Platform struct {
	Name             string `json:"name"`               // platform name
	Enabled          bool   `json:"enabled"`            // enable push to platform
	Organization     string `json:"organization"`       // set organization if you want to create repository
	Repository       string `json:"repository"`         // set repository name if you want to create repository
	URL              string `json:"url"`                // git repository url
	Username         string `json:"username"`           // platform username
	Password         string `json:"password"`           // platform password
	Token            string `json:"token"`              // platform token
	SSHKey           string `json:"ssh_key"`            // git ssh key
	SSHKeyPassphrase string `json:"ssh_key_passphrase"` // git ssh key passphrase
	RemoteName       string `json:"remote_name"`        // custom remote name
}

type Config struct {
	common.Environment
	AccessToken  string      `env:"PLUGIN_ACCESS_TOKEN" json:"PLUGIN_ACCESS_TOKEN"` // person access token
	EnvPlatforms string      `env:"PLUGIN_PLATFORMS" json:"-"`                      // platform list, string can convert to Platform slice
	Platforms    []*Platform `json:"PLUGIN_PLATFORMS,omitempty"`                    // convert from EnvPlatforms
}

func (c Config) String() string {
	c.AccessToken = common.MaskPasswordWithStars(c.AccessToken, 14)
	// Mask platform secrets
	for _, p := range c.Platforms {
		if p != nil {
			p.Token = common.MaskPasswordWithStars(p.Token, 8)
			p.Password = common.MaskPasswordWithStars(p.Password, 8)
			p.SSHKey = common.MaskPasswordWithStars(p.SSHKey, 8)
			p.SSHKeyPassphrase = common.MaskPasswordWithStars(p.SSHKeyPassphrase, 8)
		}
	}
	body, err := json.Marshal(c)
	if err != nil {
		return "show config failed: " + err.Error()
	}
	return string(body)
}

func (c *Config) Validate() error {
	if c.EnvPlatforms == "" {
		return fmt.Errorf("PLUGIN_PLATFORMS is required")
	}
	if err := json.Unmarshal([]byte(c.EnvPlatforms), &c.Platforms); err != nil {
		return fmt.Errorf("parse PLUGIN_PLATFORMS: %w", err)
	}
	if len(c.Platforms) == 0 {
		return fmt.Errorf("PLUGIN_PLATFORMS must contain at least one platform")
	}
	if c.Workspace == "" {
		return fmt.Errorf("CI_WORKSPACE is required")
	}
	for i, p := range c.Platforms {
		if !p.Enabled {
			continue
		}
		if p.URL == "" {
			return fmt.Errorf("platform[%d] (%s): url is required", i, p.Name)
		}
		if p.Token == "" && (p.Username == "" || p.Password == "") && p.SSHKey == "" {
			return fmt.Errorf("platform[%d] (%s): requires token, username+password, or ssh_key", i, p.Name)
		}
	}
	return nil
}
