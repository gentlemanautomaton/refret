package main

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gentlemanautomaton/signaler"
)

func main() {
	// Capture shutdown signals
	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)
	ctx := shutdown.Context()

	// Parse configuration
	var conf Config

	kong.Parse(&conf,
		kong.Description(description),
		kong.UsageOnError())

	fmt.Println(conf.Summary())

	// Scan the file system
	fmt.Print("Scanning directories and files...\n")
	//scanner := NewScanner(os.DirFS(conf.Root), conf.Patterns, makeShowFileCallback(conf.Matched, conf.Unmatched, conf.Verbose), conf.Concurrency)
	scanner := NewScanner(os.DirFS(conf.Root), conf.Patterns, nil, conf.Concurrency)
	scanStart := time.Now()
	files, err := scanner.Scan(ctx)
	scanEnd := time.Now()
	scanDuration := scanEnd.Sub(scanStart)
	if err != nil {
		if err == context.Canceled {
			fmt.Printf("Scanning directories and files... stopped. (%v)\n", scanDuration)
			fmt.Printf("Operation cancelled.\n")
		} else {
			fmt.Printf("Scanning directories and files... failed: %v\n", err)
		}
		os.Exit(1)
	}
	fmt.Printf("Scanning directories and files... done. (%v)\n", scanDuration)

	// Show the file scan results if requested
	if conf.Matched || conf.Unmatched {
		if err := showFiles(ctx, conf.Matched, conf.Unmatched, conf.Verbose, files); err != nil {
			if err == context.Canceled {
				fmt.Printf("Operation cancelled.\n")
			} else {
				fmt.Printf("Failed to display results: %v\n", err)
			}
			os.Exit(1)
		}
	}

	// Build the set of proposed file rename actions
	actions := BuildActions(files)
	if len(actions) == 0 {
		fmt.Printf("No actions proposed.\n")
		return
	}

	// Build the set of proposed file rename actions
	omitted := BuildOmitted(files)

	// Write the proposed actions to a TSV file
	proposedFileName := fmt.Sprintf("%s-proposed %s.tsv", conf.FileNamePrefix, currentTimestamp())
	fmt.Printf("Writing proposed actions to %s...", proposedFileName)
	err = writeActions(proposedFileName, actions)
	if err != nil {
		fmt.Printf(" failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(" done.\n")

	// Write the omitted files to a TSV file
	omittedFileName := fmt.Sprintf("%s-omitted %s.tsv", conf.FileNamePrefix, currentTimestamp())
	fmt.Printf("Writing omitted actions to %s...", proposedFileName)
	err = writeOmitted(omittedFileName, omitted)
	if err != nil {
		fmt.Printf(" failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(" done.\n")

	// Print a summary of the proposed actions
	actionsCount := pluralize(len(actions), "action", "actions")
	fmt.Printf("%s proposed.\n", actionsCount)

	// If the user hasn't opted-in to renaming things, stop now
	if !conf.Proceed {
		return
	}

	// Prompt the user for confirmation of the proposed actions
	itemsCount := pluralize(len(actions), "item", "items")
	confirmed, err := prompt(fmt.Sprintf("Proceed with rename actions affecting %s?", itemsCount))
	if err != nil {
		fmt.Printf("Cancelling due to unexpected response: %v\n", err)
		os.Exit(1)
	}

	if !confirmed {
		fmt.Printf("Cancelled.\n")
		return
	}

	// Write results to a TSV file as we make progress
	progress := make(chan Record)
	resultsFileName := fmt.Sprintf("%s-results %s.tsv", conf.FileNamePrefix, currentTimestamp())
	resultsFinished, err := writeRecordStream(resultsFileName, progress)
	if err != nil {
		fmt.Printf("Failed to prepare output file %s: %v", resultsFileName, err)
		os.Exit(1)
	}
	fmt.Printf("Progressively writing results to %s.", resultsFileName)

	// Perform the actions
	fmt.Printf("Proceeding with the proposed %s, unto whatever end.\n", actionsCount)
	processStart := time.Now()
	results, processErr := process(ctx, conf.Root, actions, progress)
	processEnd := time.Now()
	processDuration := processEnd.Sub(processStart)

	// Tell the TSV writer that we're done
	close(progress)

	// Wait for the TSV writer to finish
	writeErr := <-resultsFinished

	// Print a summary of the outcome
	{
		summary := Summarize(results)
		if processErr != nil {
			fmt.Printf("%s Stopped after %v due to error: %v\n", summary, processDuration, processErr)
			os.Exit(1)
		} else {
			fmt.Printf("%s Done. (%v)\n", summary, processDuration)
		}
	}

	// If we failed to write the results to a file, try dumping them to the screen
	if writeErr != nil {
		fmt.Printf("Failed to write results to file. Writing results to console as a last resort.\n")
		showResults(ctx, results)
	}
}
