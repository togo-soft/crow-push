package main

import "codefloe.com/actions/common"

func main() {
	var config = &Config{}
	common.ParseEnv(config)
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
