package main

// Match describes the result of a potential pattern match
type Match int

// Pattern matching results
const (
	NoPattern  Match = 0
	NotMatched Match = 1
	Matched    Match = 2
)
