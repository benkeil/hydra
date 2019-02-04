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

	dockerconfig "github.com/docker/cli/cli/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type pushCmd struct {
	version   string
	out       io.Writer
	workdir   string
	imageUtil ImageUtil
}

// PushResponse contains the response from the docker client for "docker build"
type PushResponse struct {
	Status string
	ID     string
	Error  string
}

func newPushCmd(out io.Writer, workdir string) *cobra.Command {
	c := &pushCmd{out: out, workdir: workdir, imageUtil: NewDefaultImageUtil()}

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

	// Check for Docker auth config
	dockerConfigFile, dockerErr := dockerconfig.Load("")
	check(dockerErr)

	logger.Debugf("Docker Config file: %v\n", dockerConfigFile)

	// TODO: fix the behavior with multiple YAML objects/images in config file
	registryHostname := c.imageUtil.getRegistryHostnames(config)[0]

	logger.Debugf("Registry Hostname: %v\n", registryHostname)

	registryAuth, dockerErr := dockerConfigFile.GetAuthConfig(registryHostname)
	check(dockerErr)

	registryAuthString := c.imageUtil.encodeAuth(&registryAuth)

	if registryAuthString == "" {
		registryAuthString = "hydra"
	}
	logger.Debugf("Registry Auth String: %v\n", registryAuthString)

	// Build each version from the config
	for _, v := range config.Versions {
		var directory = path.Join(c.workdir, v.Directory) + string(filepath.Separator)
		fmt.Printf("pushing %s\n", directory)
		logger.Infof("Push images from %v with tags: %v\n", directory, v.Tags)
		tags := parser.parseTags(v)
		logger.Infof("Push image tags: %v", tags)

		// Build image
		logger.Debugf("response from docker daemon:")
		for _, image := range c.imageUtil.getImageTags(config, tags) {
			fmt.Printf("push %s\n", image)
			pushResponse, err := cli.ImagePush(context.Background(), image, types.ImagePushOptions{
				RegistryAuth: registryAuthString})
			defer pushResponse.Close()
			if err != nil {
				fmt.Printf("Can not push image %s\n%s\n", image, err.Error())
			} else {
				response, err := ioutil.ReadAll(pushResponse)
				check(err)

				for _, line := range strings.Split(string(response), "\n") {
					logger.Debugf("response line: %v\n", line)
					if line != "" {
						// parse response line by line
						output := PushResponse{}
						err := json.Unmarshal([]byte(line), &output)
						check(err)

						// check the response for success/error
						// we have status for layer id
						if output.Status != "" && output.ID != "" {
							// we have status and id
							fmt.Printf("%s: %s\n", output.ID, output.Status)
						} else if output.Status != "" {
							// we have only status
							fmt.Printf("%s\n", output.Status)
						} else if output.Error != "" {
							// we have an error
							fmt.Printf("Error: %s\n", output.Error)
						}
					}
				}
			}
		}
	}

	return nil
}
