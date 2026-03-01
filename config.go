package main

import (
	"encoding/json"
	"fmt"

	"codefloe.com/actions/common"
	"github.com/tidwall/pretty"
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
	AccessToken string      `env:"PLUGIN_ACCESS_TOKEN" json:"PLUGIN_ACCESS_TOKEN"`  // person access token
	Platforms   []*Platform `env:"PLUGIN_PLATFORMS,object" json:"PLUGIN_PLATFORMS"` // platform list
}

func (c *Config) String() string {
	// Create a copy to avoid modifying original
	cfgCopy := *c
	cfgCopy.AccessToken = common.MaskPasswordWithStars(cfgCopy.AccessToken, 14)

	// Deep copy platforms
	cfgCopy.Platforms = make([]*Platform, len(c.Platforms))
	for i, p := range c.Platforms {
		if p != nil {
			pCopy := *p
			pCopy.Token = common.MaskPasswordWithStars(pCopy.Token, 8)
			pCopy.Password = common.MaskPasswordWithStars(pCopy.Password, 8)
			pCopy.SSHKey = common.MaskPasswordWithStars(pCopy.SSHKey, 8)
			pCopy.SSHKeyPassphrase = common.MaskPasswordWithStars(pCopy.SSHKeyPassphrase, 8)
			cfgCopy.Platforms[i] = &pCopy
		}
	}

	body, err := json.Marshal(cfgCopy)
	if err != nil {
		return "show config failed: " + err.Error()
	}
	prettyEnv := pretty.Pretty(body)
	return string(prettyEnv)
}

func (c *Config) Validate() error {
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
