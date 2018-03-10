package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type pushCmd struct {
	version     string
	out         io.Writer
	workdir     string
	imageHelper ImageHelper
}

// PushResponse contains the response from the docker client for "docker build"
type PushResponse struct {
	Status string
	ID     string
}

func newPushCmd(out io.Writer, workdir string) *cobra.Command {
	c := &pushCmd{out: out, workdir: workdir, imageHelper: NewDefaultImageHelper()}

	cmd := &cobra.Command{
		Use:              "push VERSION",
		Short:            "pushes docker images",
		TraverseChildren: true,
		Args:             SemverValidator(),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.version = args[0]
			return c.run()
		},
	}

	return cmd
}

func (c *pushCmd) run() error {
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
		fmt.Printf("pushing %s\n", directory)
		logger.Infof("Push images from %v with tags: %v\n", directory, v.Tags)
		tags := parser.parseTags(v)
		logger.Infof("Push image tags: %v", tags)

		// Build image
		logger.Debugf("response from docker daemon:")
		for _, image := range c.imageHelper.getImageTags(config, tags) {
			fmt.Printf("push %s\n", image)
			pushResponse, err := cli.ImagePush(context.Background(), image, types.ImagePushOptions{
				RegistryAuth: "hydra",
			})
			if err != nil {
				fmt.Printf("Can not push image %s\n%s\n", image, err.Error())
			} else {
				defer pushResponse.Close()
				response, err := ioutil.ReadAll(pushResponse)
				check(err)
				// Print response from docker daemon
				logger.Debugf("response from docker daemon:")
				for _, line := range strings.Split(string(response), "\n") {
					output := PushResponse{}
					json.Unmarshal([]byte(line), &output)
					if output.Status != "" && output.ID != "" {
						fmt.Printf("%s: %s\n", output.ID, output.Status)
					} else if output.Status != "" {
						fmt.Printf("%s\n", output.Status)
					}
				}
			}
		}
	}

	return nil
}
