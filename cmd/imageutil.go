package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/docker/cli/cli/config/types"
)

type authHeader struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ServerAddress string `json:"serveraddress"`
}

// ImageUtil provides some helper methods
type ImageUtil interface {
	getImageTags(config Config, tags []string) []string
	getRegistryHostnames(config Config) []string
	encodeAuth(authConfig *types.AuthConfig) string
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

// getRegistryHostnames returns the registry hostname:port component of image name
func (helper *DefaultImageUtil) getRegistryHostnames(config Config) []string {

	images := []string{}
	for _, image := range config.Image {
		images = append(images, strings.Split(image, "/")[0])
	}

	return images
}

// encodeAuth returns the base64 endoded JSON object for AuthConfig struct
func (helper *DefaultImageUtil) encodeAuth(authConfig *types.AuthConfig) string {
	if authConfig.Username == "" && authConfig.Password == "" {
		return ""
	}

	auth := &authHeader{
		Username:      authConfig.Username,
		Password:      authConfig.Password,
		ServerAddress: authConfig.ServerAddress,
	}
	authStr, _ := json.Marshal(auth)
	msg := []byte(authStr)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encoded, msg)
	return string(encoded)
}
