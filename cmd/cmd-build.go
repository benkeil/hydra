package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type buildCmd struct {
	version     string
	out         io.Writer
	workdir     string
	imageHelper ImageHelper
	push        bool
}

// BuildResponse contains the response from the docker client for "docker build"
type BuildResponse struct {
	Stream string
}

func newBuildCmd(out io.Writer, workdir string) *cobra.Command {
	c := &buildCmd{out: out, workdir: workdir, imageHelper: NewDefaultImageHelper()}

	cmd := &cobra.Command{
		Use:              "build VERSION",
		Short:            "builds docker images",
		TraverseChildren: true,
		Args:             SemverValidator(),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.version = args[0]
			return c.run()
		},
	}
	cmd.Flags().BoolVarP(&c.push, "push", "", false, "also push the images (experimental)")

	return cmd
}

func (c *buildCmd) run() error {
	fmt.Fprintf(c.out, "build images from workdir %s with version %s\n", c.workdir, c.version)

	// Read the config
	configReader := NewConfigReader()
	config := configReader.getConfig(c.workdir)
	logger.Tracef("Config: %+v\n", config)

	// Creates a new tag parser
	parser := NewParser(c.version)

	// Create new docker client
	cli, err := client.NewEnvClient()
	check(err)

	// Build each version from the config
	for _, v := range config.Versions {
		var directory = path.Join(c.workdir, v.Directory) + string(filepath.Separator)
		fmt.Fprintf(c.out, "building %s\n", directory)
		logger.Infof("Build images from %v with tags: %v\n", directory, v.Tags)
		tags := parser.parseTags(v)
		logger.Infof("Use tags for image: %v", tags)

		// Get Dockerfile from config or default one
		var dockerfile string
		if v.Dockerfile == "" {
			dockerfile = "Dockerfile"
		} else {
			dockerfile = v.Dockerfile
		}

		// Tar the working directory to send to the docker API
		tarfileName, err := TarWorkdir(directory)
		check(err)
		tarfile, err := os.Open(tarfileName)
		check(err)
		defer os.Remove(tarfile.Name())
		logger.Infof("Using %s for build context", tarfile.Name())

		// Build image
		buildResponse, err := cli.ImageBuild(context.Background(), tarfile, types.ImageBuildOptions{
			PullParent: true,
			Tags:       c.imageHelper.getImageTags(config, tags),
			Context:    tarfile,
			Dockerfile: dockerfile,
		})
		check(err)
		defer buildResponse.Body.Close()
		response, err := ioutil.ReadAll(buildResponse.Body)
		check(err)

		// Print response from docker daemon
		logger.Debugf("response from docker daemon:")
		for _, line := range strings.Split(string(response), "\n") {
			output := BuildResponse{}
			json.Unmarshal([]byte(line), &output)
			if output.Stream != "" {
				fmt.Printf("%s\n", strings.TrimSpace(output.Stream))
			}
		}
	}

	return nil
}
