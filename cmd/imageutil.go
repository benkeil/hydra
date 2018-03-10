package main

import (
	"fmt"
	"os"
)

// ImageUtil provides some helper methods
type ImageUtil interface {
	getImageTags(config Config, tags []string) []string
}

// DefaultImageUtil is the default implementation for the ImageUtil interface
type DefaultImageUtil struct {
}

// NewDefaultImageUtil creates a new instance of DefaultImageUtil
func NewDefaultImageUtil() *DefaultImageUtil {
	return new(DefaultImageUtil)
}

// getImageTags returns the complete image name (including registry url, image name and tag)
func (helper *DefaultImageUtil) getImageTags(config Config, tags []string) []string {
	// Define all images we want to build
	imageTags := []string{}
	for _, image := range config.Image {
		for _, tag := range tags {
			imageTags = append(imageTags, fmt.Sprintf("%s:%s", os.ExpandEnv(image), tag))
		}
	}
	return imageTags
}
