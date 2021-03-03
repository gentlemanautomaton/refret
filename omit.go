package main

import "path"

// Omit stores information about scanned files that will be omitted from the
// action list
type Omit struct {
	Path string
}

// String returns a string representation of the omitted file.
func (o Omit) String() string {
	return o.Path
}

// BuildOmitted prepares a set of file paths that have been scanned but will
// not have action taken on them.
func BuildOmitted(files []File) (omitted []Omit) {
	filter := func(file File) bool {
		return true
	}
	action := func(file File) {
		if file.Result != Matched {
			omitted = append(omitted, Omit{
				Path: path.Join(file.Parent, file.Name),
			})
		}
	}
	walkDescending(files, filter, action)
	return omitted
}
