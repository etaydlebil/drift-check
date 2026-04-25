package output

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

// Pager wraps an io.Writer and optionally pipes output through a terminal pager
// (e.g. less) when the output exceeds a threshold line count.
type Pager struct {
	w        io.Writer
	enabled  bool
	pagerCmd string
}

// NewPager returns a Pager that writes to w. When enabled is true and the
// PAGER environment variable (or the default "less") is available, long output
// is piped through it.
func NewPager(w io.Writer, enabled bool) *Pager {
	cmd := os.Getenv("PAGER")
	if cmd == "" {
		cmd = "less"
	}
	return &Pager{w: w, enabled: enabled, pagerCmd: cmd}
}

// Write implements io.Writer. Output is forwarded directly to the underlying
// writer; use WriteAll to trigger pager behaviour.
func (p *Pager) Write(b []byte) (int, error) {
	return p.w.Write(b)
}

// WriteAll writes content to the underlying writer, routing through the pager
// when enabled and the content contains more than threshold lines.
func (p *Pager) WriteAll(content string, threshold int) error {
	lines := strings.Count(content, "\n")
	if p.enabled && lines > threshold && p.pagerAvailable() {
		return p.pipeThroughPager(content)
	}
	_, err := io.WriteString(p.w, content)
	return err
}

func (p *Pager) pagerAvailable() bool {
	_, err := exec.LookPath(p.pagerCmd)
	return err == nil
}

func (p *Pager) pipeThroughPager(content string) error {
	cmd := exec.Command(p.pagerCmd) //nolint:gosec
	cmd.Stdout = p.w
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		// Fall back to direct write if we cannot obtain the pipe.
		_, werr := io.WriteString(p.w, content)
		return werr
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		_, werr := io.WriteString(p.w, content)
		return werr
	}

	_, writeErr := io.WriteString(stdin, content)
	stdin.Close()
	_ = cmd.Wait()
	return writeErr
}
