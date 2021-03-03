package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Pattern is a file matching pattern than optionally can also specify
// substitutions.
type Pattern struct {
	Expression  *regexp.Regexp
	Subtitution string
}

// UnmarshalText unmarshals the given text as a patter in p.
func (p *Pattern) UnmarshalText(text []byte) error {
	re := string(text)

	// Interpret special values as no-ops
	if re == "" || re == "_" {
		p.Expression = nil
		p.Subtitution = ""
		return nil
	}

	// If the expression contains a forward slash, interpret it as a delimiter
	// between the expression and the substitution
	sub := ""
	if pair := strings.SplitN(re, "/", 2); len(pair) == 2 {
		re = pair[0]
		sub = pair[1]
		if sub == "" {
			return errors.New("empty substitution provided in pattern")
		}
	}
	exp, err := compileRegex(re)
	if err != nil {
		return err
	}
	p.Expression = exp
	p.Subtitution = sub
	return nil
}

// String returns a string representation of the pattern.
func (p Pattern) String() string {
	if p.Expression == nil {
		return "*"
	}
	if p.Subtitution == "" {
		return fmt.Sprintf("%s", p.Expression)
	}
	return fmt.Sprintf("%s / %s", p.Expression, p.Subtitution)
}

// ApplyPattern selects an appropriate pattern for the given traversal depth,
// applies it to value, returned the result of the match.
//
// If the selected pattern supplies a substitution, it is applied and the
// substituted value is returned. Otherwise, the original value is returned.
func ApplyPattern(patterns []Pattern, depth int, value string) (result Match, newValue string) {
	if depth >= len(patterns) {
		return NoPattern, value
	}
	pattern := patterns[depth]
	if pattern.Expression == nil {
		if pattern.Subtitution != "" {
			panic("unexpected substitution with nil pattern expression")
		}
		return Matched, value
	}
	if pattern.Expression.MatchString(value) {
		if pattern.Subtitution == "" {
			return Matched, value
		}
		return Matched, pattern.Expression.ReplaceAllString(value, pattern.Subtitution)
	}
	return NotMatched, value
}

func compileRegex(re string) (*regexp.Regexp, error) {
	if re == "" {
		return nil, nil
	}

	// Force case-insensitive matching
	if !strings.HasPrefix(re, "(?i)") {
		re = "(?i)" + re
	}

	c, err := regexp.Compile(re)
	if err != nil {
		return nil, fmt.Errorf("unable to compile regular expression \"%s\": %v", re, err)
	}
	return c, nil
}
