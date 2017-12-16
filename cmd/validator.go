package main

import (
	"errors"
	"fmt"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

// SemverValidator validate if a string is a valid semver version
func SemverValidator() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("the version you want to push is required")
		}
		_, err := semver.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid semantic version: %s", args[0])
		}
		return nil
	}
}
