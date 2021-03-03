package main

import "time"

func currentTimestamp() string {
	now := time.Now()
	return now.Format("2006-01-02 150405")
}
