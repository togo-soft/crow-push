package main

import (
	"os"
	"strings"

	"codefloe.com/actions/common"
)

func main() {
	envs := os.Environ()

	for _, env := range envs {
		// 分割 key 和 value
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			common.Info("%s=%s\n", pair[0], pair[1])
		}
	}
	common.Info("============================")

	var config = &Config{}
	common.ParseEnvironment(config)
	common.Info("parse environment variables: %s", config.String())

	// validate config
	if err := config.Validate(); err != nil {
		common.Fatal("config validation: %v", err)
	}

	// run push operation
	if err := Run(config); err != nil {
		common.Fatal("push failed: %v", err)
	}
	common.Info("all platforms synced successfully")
}
