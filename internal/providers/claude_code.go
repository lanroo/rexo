package providers

import (
	"context"
	"time"
)

// ClaudeCode generates text through the Claude Code CLI (`claude`). It runs in
// non-interactive print mode and reads the prompt from stdin, which avoids
// argument-length and shell-escaping issues with long prompts.
type ClaudeCode struct {
	cli cliRunner
}

func NewClaudeCode() *ClaudeCode {
	return &ClaudeCode{cli: cliRunner{bin: "claude"}}
}

func (c *ClaudeCode) ID() string { return "claude-code" }

func (c *ClaudeCode) Detect() bool { return c.cli.available() }

func (c *ClaudeCode) Generate(ctx context.Context, req Request) (*Result, error) {
	start := time.Now()
	// `claude -p` with no inline prompt reads the prompt from stdin.
	out, err := c.cli.run(ctx, []string{"-p"}, foldSystem(req.System, req.Prompt))
	if err != nil {
		return nil, err
	}
	return &Result{
		Text:       out,
		ProviderID: c.ID(),
		Duration:   time.Since(start),
	}, nil
}
