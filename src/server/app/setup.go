// Package app performs basic app functions like setup, build time variables, loading config, registering controllers, starting http server and global definitions
package app

import (
	"fmt"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

// Setup will perform app setup or panic if it fails
func Setup() {
	// Logger setup
	{
		var err error

		if Build == BuildDev {
			err = devLoggerSetup()
			if err != nil {
				fmt.Printf("err: %v\n", err)
				panic("failed to setup logger")
			}
		} else {
			err = prodLoggerSetup()
			if err != nil {
				fmt.Printf("err: %v\n", err)
				panic("failed to setup logger")
			}
		}
	}

	// Print build time variables
	{
		zap.S().Infof("Build:      \t\t %s", Build)
		zap.S().Infof("Version:    \t\t %s", Version)
		zap.S().Infof("Commit Hash:\t\t %s", CommitHash)
		zap.S().Infof("Build Time Stamp:\t %s", BuildTimestamp)
		zap.S().Sync()
	}

	LoadConfig()

	// Dig setup
	{
		digContainer = dig.New()
	}

	// Database setup (MongoDB)
	{
		Provide(newMongoDb)
	}
}
