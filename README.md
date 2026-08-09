# refret

A regular expression file rename tool and library written in Go.

## Usage

```
Usage: refret.exe --name="migration" --root=STRING [<pattern> ...] [flags]

Searches for and optionally renames files according to regular expression
patterns. It matches file and directory names as it traverses a file system from
a given root. Successive patterns match successive traversal depths.

Proposed rename actions, omitted (non-matching) files and the results of actions
taken are logged for inspection and review.

During evaluation, files are scanned concurrently for speed. Rename operations
happen in series for safety.

Arguments:
  [<pattern> ...]    Regular expression patterns to match, with optional
                     substitution delimited by a forward slash (exp/sub)
                     ($PATTERN).

Flags:
  -h, --help                Show context-sensitive help.
      --name="migration"    Output file name prefix ($NAME).
      --root=STRING         Root path of the file directory structure ($ROOT).
  -v, --verbose             Provide verbose output ($VERBOSE).
  -m, --matched             Show matching files and directories ($MATCHED).
  -u, --unmatched           Show non-matching files and directories
                            ($UNMATCHED).
  -c, --concurrency=32      Maximum number of concurrent read operations during
                            scanning ($CONCURRENCY).
      --proceed             Proceed with renaming operations ($PROCEED).
```

## Patterns

Regular expressions are interpreted according to Go's regular expression
engine. Go has excellent documentation for the regular expression syntax
that it supports, which can be found here:

https://pkg.go.dev/regexp/syntax

The sole exception to these parsing rules is the special interpretation
of the first `/` character within a pattern, which is used by `refret` to
distinguish between the matching expression and the substitution
expression.

Substitutions can refer to capture groups by their numbered order, like `$1`
and `$2`. In the following example pattern, the first two words in each file
name would be reversed:

```
^(\w+) (\w+)/$2 $1
```
