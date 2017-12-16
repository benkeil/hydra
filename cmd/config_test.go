package main

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
)

const response = `image: my.private.registry:5000/docker-common/nginx-base
versions:
- directory: .
  args:
  dockerfile:
  tags:
    - semver
    - latest`

func TestGetConfig(t *testing.T) {
	fmt.Println("start test")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfig := NewMockConfigReader(mockCtrl)
	mockConfig.EXPECT().readConfig("some/path").Return([]byte(response), nil).Times(1)
	mockConfig.readConfig("some/path")
}
