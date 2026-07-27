package providers

import (
	"context"
	"time"
)

// Codex generates text through the OpenAI Codex CLI (`codex`). It uses the
// non-interactive `codex exec` subcommand.
type Codex struct {
	cli cliRunner
}

func NewCodex() *Codex {
	return &Codex{cli: cliRunner{bin: "codex"}}
}

func (c *Codex) ID() string { return "codex" }

func (c *Codex) Detect() bool { return c.cli.available() }

func (c *Codex) Generate(ctx context.Context, req Request) (*Result, error) {
	start := time.Now()
	out, err := c.cli.run(ctx, []string{"exec", foldSystem(req.System, req.Prompt)}, "")
	if err != nil {
		return nil, err
	}
	return &Result{
		Text:       out,
		ProviderID: c.ID(),
		Duration:   time.Since(start),
	}, nil
}
