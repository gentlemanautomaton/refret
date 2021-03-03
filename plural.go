package main

import "strconv"

func pluralize(count int, singular, plural string) string {
	switch count {
	case 1:
		return "1 " + singular
	default:
		return strconv.Itoa(count) + " " + plural
	}
}
