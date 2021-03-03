package main

import (
	"fmt"
)

const description = "Searches for and optionally renames files according to regular expression patterns. " +
	"It matches file and directory names as it traverses a file system from a given root. " +
	"Successive patterns match successive traversal depths.\n\n" +
	"Proposed rename actions, omitted (non-matching) files and the results of actions taken are logged for inspection and review.\n\n" +
	"During evaluation, files are scanned concurrently for speed. Rename operations happen in series for safety."

// Config holds configuration values ingested from the environment
// and command line.
type Config struct {
	FileNamePrefix string    `kong:"env='NAME',name='name',default='migration',required,help='Output file name prefix.'"`
	Root           string    `kong:"env='ROOT',name='root',required,help='Root path of the file directory structure.'"`
	Patterns       []Pattern `kong:"env='PATTERN',name='pattern',arg,optional,help='Regular expression patterns to match, with optional substitution delimited by a forward slash (exp/sub).'"`
	Verbose        bool      `kong:"env='VERBOSE',name='verbose',short='v',help='Provide verbose output.'"`
	Matched        bool      `kong:"env='MATCHED',name='matched',short='m',help='Show matching files and directories.'"`
	Unmatched      bool      `kong:"env='UNMATCHED',name='unmatched',short='u',help='Show non-matching files and directories.'"`
	Concurrency    int       `kong:"env='CONCURRENCY',name='concurrency',short='c',default='32',help='Maximum number of concurrent read operations during scanning.'"`
	Proceed        bool      `kong:"env='PROCEED',name='proceed',help='Proceed with renaming operations.'"`
}

// Summary returns a multiline string describing the configuration.
func (conf Config) Summary() string {
	output := fmt.Sprintf("Output File Name Prefix: %s", conf.FileNamePrefix)
	output += fmt.Sprintf("\nBase Path (Root): %s", conf.Root)
	for depth, pattern := range conf.Patterns {
		if pattern.Expression == nil {
			output += fmt.Sprintf("\nDepth %d Pattern: .*", depth)
		} else {
			output += fmt.Sprintf("\nDepth %d Pattern: %s", depth, pattern.Expression)
		}
		if pattern.Subtitution != "" {
			output += fmt.Sprintf("\nDepth %d Substitution: %s", depth, pattern.Subtitution)
		}
	}
	switch {
	case conf.Matched && conf.Unmatched:
		output += fmt.Sprintf("\nShow: Both matching and non-matching files and directories")
	case conf.Matched:
		output += fmt.Sprintf("\nShow: Matching files and directories")
	case conf.Unmatched:
		output += fmt.Sprintf("\nShow: Non-matching files and directories")
	}
	if conf.Verbose {
		output += fmt.Sprintf("\nVerbose Output")
	}
	output += fmt.Sprintf("\nConcurrency: %d", conf.Concurrency)
	if conf.Proceed {
		output += fmt.Sprintf("\nExecution Requested")
	}
	return output
}
