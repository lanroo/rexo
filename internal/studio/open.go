package studio

import (
	"os/exec"
	"runtime"
)

// openBrowser tries to open url in the default browser. Any failure is
// non-fatal — the URL is always printed so the user can open it manually.
func openBrowser(url string) {
	var name string
	var args []string
	switch runtime.GOOS {
	case "windows":
		name, args = "cmd", []string{"/c", "start", "", url}
	case "darwin":
		name, args = "open", []string{url}
	default:
		name, args = "xdg-open", []string{url}
	}
	_ = exec.Command(name, args...).Start()
}
