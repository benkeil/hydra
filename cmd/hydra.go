package main

import (
	"fmt"
	"os"

	"github.com/juju/loggo"
	"github.com/spf13/cobra"
)

// http://how-bazaar.blogspot.de/2013/10/loggo-hierarchical-loggers-for-go.html
var logger = loggo.GetLogger("hydra")

var version string

func init() {
	loggo.ConfigureLoggers(fmt.Sprintf("<root>=INFO; hydra=%s", "ERROR"))
}

func main() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	var workdir string
	var debug bool
	cmd := &cobra.Command{
		Use:          "hydra",
		Short:        "Hydra builds docker images and add multiple convenient tags",
		Version:      version,
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var logLevel = "ERROR"
			if debug {
				logLevel = "DEBUG"
			}
			loggo.ConfigureLoggers(fmt.Sprintf("<root>=INFO; hydra=%s", logLevel))
		},
	}

	cmd.PersistentFlags().StringVarP(&workdir, "workdir", "w", ".", "the root directory of the project you want to build")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output")
	cmd.PersistentFlags().Parse(args)

	out := cmd.OutOrStdout()

	cmd.AddCommand(
		newBuildCmd(out, workdir),
		newPushCmd(out, workdir),
	)

	return cmd
}

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
		logger.Errorf(e.Error())
		panic(e)
	}
}
