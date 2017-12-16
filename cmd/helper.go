package main

import (
	"fmt"
)

// ImageHelper
type ImageHelper interface {
	getImageTags(config Config, tags []string) []string
}

// DefaultImageHelper is the default implementation for the ImageHelper interface
type DefaultImageHelper struct {
}

// NewDefaultImageHelper creates a new instance of DefaultImageHelper
func NewDefaultImageHelper() *DefaultImageHelper {
	return new(DefaultImageHelper)
}

func (helper *DefaultImageHelper) getImageTags(config Config, tags []string) []string {
	// Define all images we want to build
	imageTags := []string{}
	for _, image := range config.Image {
		for _, tag := range tags {
			imageTags = append(imageTags, fmt.Sprintf("%s:%s", image, tag))
		}
	}
	return imageTags
}
