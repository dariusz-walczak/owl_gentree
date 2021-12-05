package main

import (
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
)

type AppArgs struct {
	LogLevel log.Level
}

// Parse the command line arguments and return the results
// The returned structure is parser agnostic
func parseArgs() (AppArgs, error) {
	var def struct {
		LogLevel string `long:"log-level" choice:"trace" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal" choice:"panic" default:"info"`
	}

	_, err := flags.Parse(&def)
	if err != nil {
		if fErr, ok := err.(*flags.Error); ok {
			if fErr.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}

		return AppArgs{}, err
	}

	// Convert the log level argument value (string) to the logrus log level:
	level, err := log.ParseLevel(def.LogLevel)
	if err != nil {
		return AppArgs{}, err
	}

	// Return the final args structure:
	return AppArgs{
		LogLevel: level,
	}, nil
}
