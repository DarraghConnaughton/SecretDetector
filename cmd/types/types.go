package types

import "regexp"

type Rule struct {
	ID          string   `toml:"id"`
	Description string   `toml:"description"`
	Regex       string   `toml:"regex"`
	Tags        []string `toml:"tags"`
	Keywords    []string `toml:"keywords"`
}

// Secrets can be expanded to include the line number as well.
type Report struct {
	Timestamp int64  `json:"timestamp"`
	Secrets   []Line `json:"secrets"`
}

type Line struct {
	Number   int
	Contents string
	Filename string
}

type TomlConfig struct {
	Rules []Rule `toml:"rules"`
}

type Context struct {
	FilePaths      []string
	SecretPatterns []*regexp.Regexp
	PatternConfig  string
	StartTime      int64
	EndTime        int64
}
