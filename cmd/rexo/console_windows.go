//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

// launchedByDoubleClick reports whether the process owns its console alone,
// which on Windows means it was started from Explorer (a double-click) rather
// than from an existing terminal. When launched from cmd/PowerShell, the shell
// is also attached to the console, so the process count is 2 or more.
func launchedByDoubleClick() bool {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleProcessList")
	var pids [2]uint32
	n, _, _ := proc.Call(uintptr(unsafe.Pointer(&pids[0])), uintptr(len(pids)))
	return n == 1
}
