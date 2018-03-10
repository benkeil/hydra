package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/juju/loggo"
	"github.com/spf13/cobra"
)

// http://how-bazaar.blogspot.de/2013/10/loggo-hierarchical-loggers-for-go.html
var logger = loggo.GetLogger("hydra")

var version string

func init() {
	loggo.ConfigureLoggers("<root>=INFO; hydra=ERROR")
}

func main() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func emptyRun(*cobra.Command, []string) {}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func newRootCmd(args []string) *cobra.Command {
	var workdir string
	cmd := &cobra.Command{
		Use:          "hydra",
		Short:        "Hydra builds docker images and add multiple convenient tags",
		Version:      version,
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringVarP(&workdir, "workdir", "w", ".", "the root directory of the project you want to build")
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
