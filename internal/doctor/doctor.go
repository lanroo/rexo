package doctor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/lanroo/rexo/internal/providers"
)

type Check struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Details string `json:"details"`
}

type Report struct {
	Healthy bool    `json:"healthy"`
	OS      string  `json:"os"`
	Arch    string  `json:"arch"`
	Checks  []Check `json:"checks"`
}

func Run() Report {
	report := Report{
		Healthy: true,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	}
	report.Checks = append(report.Checks,
		platformCheck(),
		executableCheck("git"),
		configDirectoryCheck(),
		aiProvidersCheck(),
	)
	for _, check := range report.Checks {
		if !check.OK {
			report.Healthy = false
		}
	}
	return report
}

func Print(w io.Writer, report Report) {
	fmt.Fprintf(w, "REXO doctor (%s/%s)\n", report.OS, report.Arch)
	for _, check := range report.Checks {
		status := "PASS"
		if !check.OK {
			status = "FAIL"
		}
		fmt.Fprintf(w, "[%s] %s: %s\n", status, check.Name, check.Details)
	}
	if report.Healthy {
		fmt.Fprintln(w, "Result: healthy")
	} else {
		fmt.Fprintln(w, "Result: action required")
	}
}

func platformCheck() Check {
	supportedOS := runtime.GOOS == "windows" || runtime.GOOS == "darwin" || runtime.GOOS == "linux"
	supportedArch := runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64"
	return Check{
		Name:    "platform",
		OK:      supportedOS && supportedArch,
		Details: runtime.GOOS + "/" + runtime.GOARCH,
	}
}

func executableCheck(name string) Check {
	path, err := exec.LookPath(name)
	if err != nil {
		return Check{Name: name, OK: false, Details: "not found in PATH"}
	}
	return Check{Name: name, OK: true, Details: path}
}

func configDirectoryCheck() Check {
	path, err := os.UserConfigDir()
	if err != nil {
		return Check{Name: "config directory", OK: false, Details: err.Error()}
	}
	return Check{Name: "config directory", OK: true, Details: path}
}

// aiProvidersCheck reports which text.generate provider CLIs are on PATH. It is
// advisory, not health-gating: the deterministic core runs without any AI
// provider, but `rexo demo` needs at least one.
func aiProvidersCheck() Check {
	available := providers.DefaultResolver("", "").Available()
	if len(available) == 0 {
		return Check{
			Name:    "ai-providers",
			OK:      true,
			Details: "none found — install claude (Claude Code), codex, or ollama to use `rexo demo`",
		}
	}
	return Check{Name: "ai-providers", OK: true, Details: strings.Join(available, ", ")}
}
