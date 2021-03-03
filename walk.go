package main

// FileFilter is a function that returns true for files that should be
// included in a result set.
type FileFilter func(File) bool

// FileAction is a function that takes action on a file.
type FileAction func(File)

// walkDescending visits every file in descending depth-first traversal and
// performs the requested action. It starts at the root and works its way
// its down to the leaves.
//
// If filter is non-nil, it is called before descending into a directory
// or taking action on a file.
func walkDescending(files []File, filter FileFilter, action FileAction) {
	for _, file := range files {
		if filter == nil || filter(file) {
			action(file)
			walkDescending(file.Contents, filter, action)
		}
	}
}

// walkAscending visits every file in ascending depth-first traversal and
// performs the requested action. It starts at the leaves and works
// its way up to the root.
//
// If filter is non-nil, it is called before descending into a directory
// or taking action on a file.
func walkAscending(files []File, filter func(File) bool, action func(File)) {
	for _, file := range files {
		if filter == nil || filter(file) {
			walkAscending(file.Contents, filter, action)
			action(file)
		}
	}
}
