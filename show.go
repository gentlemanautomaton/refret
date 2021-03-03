package main

import (
	"context"
	"fmt"
	"strings"
)

// showFiles prints the results of a file scan.
func showFiles(ctx context.Context, matched, unmatched, verbose bool, files []File) error {
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return err
		}

		if shouldInclude(file, matched, unmatched) {
			if verbose {
				fmt.Printf("%s%s\n", strings.Repeat("  ", file.Depth), file.VerboseString())
			} else {
				fmt.Printf("%s%s\n", strings.Repeat("  ", file.Depth), file)
			}
		}

		if err := showFiles(ctx, matched, unmatched, verbose, file.Contents); err != nil {
			return err
		}
	}
	return nil
}

func makeShowFileCallback(matched, unmatched, verbose bool) ScanCallback {
	return func(file File) {
		if shouldInclude(file, matched, unmatched) {
			if verbose {
				fmt.Printf("%s %s\n", strings.Repeat("  ", file.Depth), file.VerboseString())
			} else {
				fmt.Printf("%s %s\n", strings.Repeat("  ", file.Depth), file)
			}
		}
	}
}

func shouldInclude(file File, matched, unmatched bool) bool {
	if matched {
		if file.Result == Matched {
			return true
		}
		if file.DescendantsMatched > 0 {
			return true
		}
	}
	if unmatched {
		if file.Result == NotMatched {
			return true
		}
		if file.DescendantsNotMatched > 0 {
			return true
		}
	}
	return false
}

func showResults(ctx context.Context, results []Record) error {
	for i, record := range results {
		if err := ctx.Err(); err != nil {
			return err
		}
		fmt.Printf("%d: %s\n", i, record)
	}
	return nil
}
