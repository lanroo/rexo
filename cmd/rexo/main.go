package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lanroo/rexo/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	code := cli.Run(os.Args[1:], os.Stdout, os.Stderr, cli.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	})

	// When double-clicked in Explorer, Windows opens a console just for us and
	// closes it the instant we exit — the window would only flash. Keep it open
	// so the person can read the welcome screen and knows nothing crashed.
	if len(os.Args) <= 1 && launchedByDoubleClick() {
		fmt.Print("\nPress Enter to close this window...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}

	os.Exit(code)
}
