package main

import (
	"fmt"
	"io/fs"
	"strconv"
)

// File represents a file within the target file system
type File struct {
	Depth                 int
	Index                 int
	Name                  string
	NewName               string
	Parent                string
	NewParent             string
	IsDir                 bool
	Result                Match
	DescendantsMatched    int
	DescendantsNotMatched int
	DescendantActions     int
	Contents              []File
}

// NewFile returns a file with it static properties set.
func NewFile(patterns []Pattern, depth int, index int, entry fs.DirEntry, oldParent, newParent string) File {
	name := entry.Name()
	result, newName := ApplyPattern(patterns, depth, name)
	return File{
		Depth:     depth,
		Index:     index,
		Name:      name,
		NewName:   newName,
		Result:    result,
		Parent:    oldParent,
		NewParent: newParent,
		IsDir:     entry.IsDir(),
	}
}

// Actionable returns true if action is needed for f.
func (f File) Actionable() bool {
	return f.NewName != f.Name
}

// VerboseString returns a more verbose string representation of f.
func (f File) VerboseString() string {
	return fmt.Sprintf("%6s [%s] [%s] %s", strconv.Itoa(f.Index)+":", f.kind(), f.result(), f.name())
}

// String returns a string representation of f.
func (f File) String() string {
	return fmt.Sprintf("[%s] %s", f.kindShort(), f.name())
}

func (f File) name() string {
	if f.NewName == "" {
		return f.Name
	}
	if f.NewName != f.Name {
		return fmt.Sprintf("%s → %s", f.Name, f.NewName)
	}
	return f.Name
}

func (f File) kind() string {
	if f.IsDir {
		return "DIR"
	}
	return "FILE"
}

func (f File) kindShort() string {
	if f.IsDir {
		return "D"
	}
	return "F"
}

func (f File) result() string {
	switch f.Result {
	case Matched:
		return "*"
	case NotMatched:
		return " "
	default:
		return "~"
	}
}

// FileIter iterates through a set of files.
type FileIter struct {
	files []File
	index int
}

// NewFileIter returns a new file iterator for files.
func NewFileIter(files []File) *FileIter {
	return &FileIter{files, 0}
}

// Next returns the index and a pointer to the next file.
//
// If the iterator has reached the end of the series, it returns -1 and nil.
func (iter *FileIter) Next() (int, *File) {
	if iter.Done() {
		return -1, nil
	}
	i, file := iter.index, &iter.files[iter.index]
	iter.index++
	return i, file
}

// Done returns true if the iterator has reached the end of the series.
func (iter *FileIter) Done() bool {
	return iter.index >= len(iter.files)
}
