package main

import "path"

// Action describes a proposed action on a file system.
type Action struct {
	OldPath string
	NewPath string
}

// String returns a string representation of the action.
func (a Action) String() string {
	return a.OldPath + " → " + a.NewPath
}

// BuildActions prepares a set of actions to be taken on a set of files that
// have been scanned. All file operations that match the requested patterns
// and are actionable will be included.
func BuildActions(files []File) (actions []Action) {
	filter := func(file File) bool {
		return file.Result == Matched
	}
	action := func(file File) {
		if file.Actionable() {
			actions = append(actions, Action{
				OldPath: path.Join(file.Parent, file.Name),
				NewPath: path.Join(file.Parent, file.NewName),
			})
		}
	}
	walkAscending(files, filter, action)
	return actions
}
