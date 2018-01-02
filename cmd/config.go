package main

//go:generate mockgen -source=config.go -destination=config_mock.go -package=main

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const hydraFile = "hydra.yaml"

// ConfigReader reads and parse the hydra.yaml
type ConfigReader interface {
	getConfig(directory string) Config
	readConfig(path string) (data []byte, err error)
	parseConfig(data []byte) Config
}

// DefaultConfigReader is the default ConfigReader implementation
type DefaultConfigReader struct{}

// Config contains basic informations about the project
type Config struct {
	Image    []string
	Versions []Version
}

// Version contains informations that are necessary to build and tag the image
type Version struct {
	Directory  string
	Tags       []string
	Args       []string
	Dockerfile string
}

// NewConfigReader returns a new ConfigReader
func NewConfigReader() ConfigReader {
	return new(DefaultConfigReader)
}

// GetConfig returns the parsed hydra config
func (c *DefaultConfigReader) getConfig(directory string) Config {
	data, err := c.readConfig(path.Join(directory, hydraFile))
	check(err)
	return c.parseConfig(data)
}

func (c *DefaultConfigReader) readConfig(path string) (data []byte, err error) {
	return ioutil.ReadFile(path)
}

func (c *DefaultConfigReader) parseConfig(data []byte) Config {
	config := Config{}
	err := yaml.Unmarshal([]byte(data), &config)
	check(err)
	return config
}
