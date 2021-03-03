package main

// Record is a migration record describing actions taken on a particular
// file or folder.
type Record struct {
	OldPath string
	NewPath string
	Error   string
}

// String returns a string representation of the record.
func (r Record) String() string {
	if r.Error != "" {
		return "FAIL: " + r.OldPath + " → " + r.NewPath + ": " + r.Error
	}
	return "MOVE: " + r.OldPath + " → " + r.NewPath
}
