package main

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
)

const semverTag string = "semver"

// TagParser is the interface for the Parser
type TagParser interface {
	parseTags(version Version) []string
}

// SemverTagParser parse semver based version
type SemverTagParser struct {
	semverVersion semver.Version
}

// parseTags parse all configured tags
func (parser *SemverTagParser) parseTags(version Version) []string {
	var parsedTags []string
	for _, tag := range version.Tags {
		logger.Debugf("Parse tag: %v\n", tag)
		parsedTags = append(parsedTags, parser.parseTag(version, tag)...)
	}
	if len(version.Tags) == 0 {
		logger.Debugf("Using DefaultStrategy for directory: %s", version.Directory)
		strategy := DefaultStrategy{version, parser.semverVersion.String()}
		parsedTags = append(parsedTags, strategy.GetTags()...)
	}
	return parsedTags
}

// parseTag parse just one tag and returns all corresponding tags
func (parser *SemverTagParser) parseTag(version Version, tag string) []string {
	var parsedTags []string
	if strings.Contains(tag, semverTag) {
		logger.Debugf("Using SemverStrategy for tag: %v", tag)
		strategy := SemverStrategy{version, parser.semverVersion, tag}
		parsedTags = append(parsedTags, strategy.GetTags()...)
	} else {
		logger.Debugf("Using SimpleStrategy for tag: %v", tag)
		strategy := SimpleStrategy{tag}
		parsedTags = append(parsedTags, strategy.GetTags()...)
	}
	logger.Debugf("Parsed tags: %v\n", parsedTags)
	return parsedTags
}

// DefaultTagParser parse not semver based versions
type DefaultTagParser struct {
	version string
}

// parseTags parse all configured tags
func (parser *DefaultTagParser) parseTags(version Version) []string {
	var parsedTags []string
	for _, tag := range version.Tags {
		logger.Debugf("Parse tag: %v\n", tag)
		parsedTags = append(parsedTags, parser.parseTag(version, tag)...)
	}
	if len(version.Tags) == 0 {
		logger.Debugf("Using DefaultStrategy for directory: %s", version.Directory)
		strategy := DefaultStrategy{version, parser.version}
		parsedTags = append(parsedTags, strategy.GetTags()...)
	}
	return parsedTags
}

// parseTag parse just one tag and returns all corresponding tags
func (parser *DefaultTagParser) parseTag(version Version, tag string) []string {
	var parsedTags []string
	if strings.Contains(tag, semverTag) {
		logger.Debugf("Using ReplaceStrategy for tag: %v", tag)
		strategy := ReplaceStrategy{tag, parser.version}
		parsedTags = append(parsedTags, strategy.GetTags()...)
	} else {
		// convenient should not be build if we dont have a semver tag
		logger.Debugf("Using SkipStrategy for tag: %v", tag)
	}
	logger.Debugf("Parsed tags: %v\n", parsedTags)
	return parsedTags
}

// TagParserFactory is a factory for a TagParser
type TagParserFactory func(version string) TagParser

// NewDefaultTagParser creates a new DefaultTagParser
func NewDefaultTagParser(version string) TagParser {
	return &DefaultTagParser{version: version}
}

// NewSemverTagParser creates a new SemverTagParser
func NewSemverTagParser(version semver.Version) TagParser {
	return &SemverTagParser{semverVersion: version}
}

// NewParser creates an new Parser
func NewParser(version string) TagParser {
	v, err := semver.Parse(version)
	if err != nil {
		return NewDefaultTagParser(version)
	}
	return NewSemverTagParser(v)
}

// TagStrategy is the interface for all tagging strategies
type TagStrategy interface {
	GetTags() []string
}

// DefaultStrategy returns the default tag for an image
type DefaultStrategy struct {
	Version Version
	version string
}

// GetTags returns the default tags for an image
func (strategy *DefaultStrategy) GetTags() []string {
	return []string{fmt.Sprintf("%s-%s", strategy.version, strings.Replace(strategy.Version.Directory, "/", "-", -1))}
}

// SemverStrategy if a tags contains the string "semver", this strategy splits it into three convenient tags
type SemverStrategy struct {
	Version       Version
	semverVersion semver.Version
	Tag           string
}

// GetTags returns the semver tags for an image
func (strategy *SemverStrategy) GetTags() []string {
	return []string{
		strings.Replace(strategy.Tag, semverTag, fmt.Sprintf("%d.%d.%d", strategy.semverVersion.Major, strategy.semverVersion.Minor, strategy.semverVersion.Patch), -1),
		strings.Replace(strategy.Tag, semverTag, fmt.Sprintf("%d.%d", strategy.semverVersion.Major, strategy.semverVersion.Minor), -1),
		strings.Replace(strategy.Tag, semverTag, fmt.Sprintf("%d", strategy.semverVersion.Major), -1),
	}
}

// SimpleStrategy just returns the tag defined in the config
type SimpleStrategy struct {
	Tag string
}

// GetTags returns a simple tag like "latest" for an image
func (strategy *SimpleStrategy) GetTags() []string {
	return []string{strategy.Tag}
}

// ReplaceStrategy just returns the tag defined in the config
type ReplaceStrategy struct {
	Tag     string
	Version string
}

// GetTags returns a simple tag like "latest" for an image
func (strategy *ReplaceStrategy) GetTags() []string {
	return []string{
		strings.Replace(strategy.Tag, semverTag, fmt.Sprintf("%s", strategy.Version), -1),
	}
}
