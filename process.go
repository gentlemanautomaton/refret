package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// process performs the given set of file rename actions and returns the
// results.
//
// If progress is non-nil, records of actions taken will be sent on the
// channel as they occur.
//
// If any error is returned, the set of completed records will be returned with it.
func process(ctx context.Context, root string, actions []Action, progress chan<- Record) (results []Record, err error) {
	for i, action := range actions {
		if err := ctx.Err(); err != nil {
			return results, err
		}

		// Combine the relative action paths with the root. Use filepath.Join
		// here, because the paths need to make sense to the local file
		// system.
		from := filepath.Join(root, action.OldPath)
		to := filepath.Join(root, action.NewPath)

		// Let the user know what we're doing
		fmt.Printf("Performing action %d: %s → %s\n", i, from, to)

		// Perform the action and print an error if it failed
		err := os.Rename(from, to)

		// Marshal the result as a record
		record := Record{OldPath: from, NewPath: to}

		// Record and print an error if the action failed
		if err != nil {
			record.Error = err.Error()
			fmt.Printf("  FAILED: %v\n", err)
		}

		// Send the record to the progress channel for logging
		if progress != nil {
			progress <- record
		}

		// Apppend the record to the result set
		results = append(results, record)
	}
	return results, nil
}
