# Hydra

[![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)
[![Build Status](https://travis-ci.org/benkeil/hydra.svg?branch=master)](https://travis-ci.org/benkeil/hydra) [![Go Report Card](https://goreportcard.com/badge/github.com/benkeil/hydra)](https://goreportcard.com/report/github.com/benkeil/hydra) [![codecov](https://codecov.io/gh/benkeil/hydra/branch/master/graph/badge.svg)](https://codecov.io/gh/benkeil/hydra) [![Github Release](https://img.shields.io/github/release/benkeil/hydra.svg)](https://github.com/benkeil/hydra/releases)

Hydra helps you to build docker images of your applications with [semver](https://semver.org) based tags.

## How it works

Hydra uses a config file named `hydra.yaml` to configure your docker images of your application and generates multiple tags like the offical docker community images (for example the [golang](https://hub.docker.com/_/golang/) image).

Here is an example `hydra.yaml` for a typical php base image in that companies use:

```yaml
image:
- my.private.registry:5000/docker-common/php-base
versions:
- directory: php5.6/alpine
  tags:
  - semver-php5.6-alpine
  - php5.6-alpine
  - semver-alpine
  - php5.6
  - latest
- directory: php5.6/debian
  tags:
  - semver-php5.6-debian
  - php5.6-debian
- directory: php7.1/debian
- directory: php7.1/alpine
  tags:
  - semver-php7.1-alpine
  - php7.1-alpine
  - php7.1
```

After you build your project with `hydra build 1.3.5` you get the following images:

```bash
build images from workdir examples/php-base/ with version 1.3.5
building examples/php-base/php5.6/alpine/
Step 1/1 : FROM php:5.6-alpine
---> cad28366b86f
Successfully built cad28366b86f
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3.5-php5.6-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3-php5.6-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:1-php5.6-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:php5.6-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:php5.6
Successfully tagged my.private.registry:5000/docker-common/php-base:latest
building examples/php-base/php5.6/debian/
Step 1/1 : FROM php:5.6-jessie
---> ee5bce1c39ee
Successfully built ee5bce1c39ee
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3.5-php5.6-debian
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3-php5.6-debian
Successfully tagged my.private.registry:5000/docker-common/php-base:1-php5.6-debian
Successfully tagged my.private.registry:5000/docker-common/php-base:php5.6-debian
building examples/php-base/php7.1/debian/
Step 1/1 : FROM php:7.1-jessie
---> 7e10b050a58c
Successfully built 7e10b050a58c
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3.5-php7.1-debian
building examples/php-base/php7.1/alpine/
Step 1/1 : FROM php:7.1-alpine
---> 07ecc747a915
Successfully built 07ecc747a915
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3.5-php7.1-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3-php7.1-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:1-php7.1-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3.5-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:1.3-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:1-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:php7.1-alpine
Successfully tagged my.private.registry:5000/docker-common/php-base:php7.1
Successfully tagged my.private.registry:5000/docker-common/php-base:alpine
```

## Tagging strategies

### Default

Will also be aplied if no tag is specified.

    {SEMVER-VERSION}-{DIRECTORY-PATH}

### Simple

Just adds the tag like in the config (e.g. `latest`).

### Semver

The string `semver` is a special tag that generates three convenient tags. It can be at any position in the string and will be replaced.

    {MAJOR-VERSION}.{FEATURE-VERSION}.{BUGFIX-VERSION}[-{SUFFIX}]
    {MAJOR-VERSION}.{FEATURE-VERSION}[-{SUFFIX}]
    {MAJOR-VERSION}[-{SUFFIX}]

## Commands

### Push

Can be used to push all images to the registry.

    hydra push VERSION -w WORKDIR
