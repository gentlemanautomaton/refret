package main

import "fmt"

// Summary holds summarized data for a set of records.
type Summary struct {
	MoveAttempted int
	MoveSuccess   int
	MoveFailure   int
}

// String returns a string representation of s.
func (s Summary) String() string {
	if s.MoveAttempted <= 0 {
		return "No moves attempted."
	}
	movesCount := pluralize(s.MoveAttempted, "move", "moves")
	switch {
	case s.MoveSuccess > 0 && s.MoveFailure > 0:
		return fmt.Sprintf("%d of %s succeeded. %d failed.", s.MoveSuccess, movesCount, s.MoveFailure)
	case s.MoveSuccess > 0:
		return fmt.Sprintf("%d of %s succeeded.", s.MoveSuccess, movesCount)
	case s.MoveFailure > 0:
		return fmt.Sprintf("%d of %s failed.", s.MoveFailure, movesCount)
	default:
		return fmt.Sprintf("%s attempted, but nothing happened.", movesCount)
	}
}

// Summarize returns a summary for a set of records.
func Summarize(records []Record) Summary {
	var s Summary
	for _, record := range records {
		s.MoveAttempted++
		if record.Error == "" {
			s.MoveSuccess++
		} else {
			s.MoveFailure++
		}
	}
	return s
}
