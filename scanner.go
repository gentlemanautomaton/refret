package main

import (
	"context"
	"fmt"
	"io/fs"
	"path"

	"golang.org/x/sync/errgroup"
)

// ScanCallback is a callback function that can be called for each file as
// it is scanned.
type ScanCallback func(file File)

// Scanner scans file systems for matching files.
type Scanner struct {
	root     fs.FS
	patterns []Pattern
	callback ScanCallback
	pool     *Pool
}

// NewScanner prepares a new scanner for the given file system and pattern.
//
// If callback is supplied, it will be called for each file scanned in order.
func NewScanner(root fs.FS, patterns []Pattern, callback ScanCallback, concurrency int) *Scanner {
	return &Scanner{
		root:     root,
		patterns: patterns,
		callback: callback,
		pool:     NewPool(concurrency),
	}
}

// Scan returns the result of scanning for files based on the given
// configuration.
func (s Scanner) Scan(ctx context.Context) (files []File, err error) {
	entries, err := fs.ReadDir(s.root, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to scan root: %v", err)
	}

	files = make([]File, len(entries))
	for i := range entries {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		files[i] = NewFile(s.patterns, 0, i, entries[i], "", "")
	}

	group, ctx := errgroup.WithContext(ctx)
	if err := s.walk(ctx, group, files); err != nil {
		group.Wait() // Make sure all of the goroutines wrap up
		return nil, err
	}

	return files, nil
}

func (s Scanner) walk(ctx context.Context, group *errgroup.Group, files []File) error {
	pass1 := NewFileIter(files)
	pass2 := NewFileIter(files)
	ready := make([]<-chan error, len(files))

	for !pass1.Done() || !pass2.Done() {
		// Pass 1 enumerates directories and stays ahead of pass 2
		added := 0
		for i, file := pass1.Next(); file != nil; i, file = pass1.Next() {
			if !s.shouldTraverse(*file) {
				continue
			}
			r, err := s.enumerate(ctx, group, file)
			if err != nil {
				return err
			}
			ready[i] = r
			added++
			if added >= 2 {
				break
			}
		}

		// Pass 2 waits for enumeration to complete, then fires the callback
		// and descends
		for i, file := pass2.Next(); file != nil; i, file = pass2.Next() {
			if ready[i] != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case err := <-ready[i]:
					if err != nil {
						return err
					}
				}
			}
			if s.callback != nil {
				s.callback(*file)
			}
			if err := s.walk(ctx, group, file.Contents); err != nil {
				return err
			}
			file.DescendantsMatched = countDescendantsMatched(file.Contents)
			file.DescendantsNotMatched = countDescendantsNotMatched(file.Contents)
			file.DescendantActions = countDescendantActions(file.Contents)
			break
		}
	}

	return nil
}

func (s Scanner) enumerate(ctx context.Context, g *errgroup.Group, file *File) (<-chan error, error) {
	if err := s.pool.Acquire(ctx); err != nil {
		return nil, err
	}
	c := make(chan error)
	go func() {
		defer s.pool.Release()
		defer close(c)

		entries, err := fs.ReadDir(s.root, path.Join(file.Parent, file.Name))
		if err != nil {
			c <- fmt.Errorf("failed to collect contents of subdirectory: %v", err)
			return
		}

		contents := make([]File, len(entries))
		for i := range entries {
			if err := ctx.Err(); err != nil {
				c <- err
				return
			}
			contents[i] = NewFile(s.patterns, file.Depth+1, i, entries[i], path.Join(file.Parent, file.Name), path.Join(file.NewParent, file.NewName))
		}

		file.Contents = contents
	}()
	return c, nil
}

func (s Scanner) shouldTraverse(file File) bool {
	return file.IsDir && file.Result == Matched && (file.Depth+1 < len(s.patterns))
}

func countDescendantsMatched(files []File) int {
	count := 0
	for i := range files {
		if files[i].Result == Matched {
			count++
		}
		count += files[i].DescendantsMatched
	}
	return count
}

func countDescendantsNotMatched(files []File) int {
	count := 0
	for i := range files {
		if files[i].Result == NotMatched {
			count++
		}
		count += files[i].DescendantsNotMatched
	}
	return count
}

func countDescendantActions(files []File) int {
	count := 0
	for i := range files {
		if files[i].Actionable() {
			count++
		}
		count += files[i].DescendantActions
	}
	return count
}
