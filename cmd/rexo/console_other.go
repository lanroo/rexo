//go:build !windows

package main

// launchedByDoubleClick is Windows-specific; on other platforms a CLI run from
// a file manager still uses the terminal that launched it, so there is nothing
// to keep open.
func launchedByDoubleClick() bool { return false }
