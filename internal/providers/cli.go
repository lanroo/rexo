package providers

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// ansiEscape matches ANSI/VT100 control sequences (CSI). Some CLIs leak cursor
// and colour codes into stdout; stripping them keeps captured text clean.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;?=]*[ -/]*[@-~]`)

// cliRunner runs a local CLI as a generation backend. Keeping the exec details
// in one place means each provider adapter only declares its binary name and
// how to build its argument list.
type cliRunner struct {
	bin string
}

func (r cliRunner) available() bool {
	_, err := exec.LookPath(r.bin)
	return err == nil
}

// run executes the binary with args, feeding stdin, and returns trimmed stdout.
// A non-zero exit includes stderr in the error so failures are diagnosable.
func (r cliRunner) run(ctx context.Context, args []string, stdin string) (string, error) {
	cmd := exec.CommandContext(ctx, r.bin, args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// CLIs report failures inconsistently: some use stderr, some print the
		// reason to stdout (e.g. "Not logged in"). Prefer stderr, fall back to
		// stdout, then to the raw exec error, so the cause is never swallowed.
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("%s failed: %s", r.bin, msg)
	}
	return strings.TrimSpace(ansiEscape.ReplaceAllString(stdout.String(), "")), nil
}

// foldSystem merges an optional system prompt into the user prompt. CLI tools
// vary in how (or whether) they accept a separate system role, so folding it in
// is the portable path and keeps the cache key stable across adapters.
func foldSystem(system, prompt string) string {
	if strings.TrimSpace(system) == "" {
		return prompt
	}
	return system + "\n\n" + prompt
}
