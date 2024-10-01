package ssh

import (
	"bufio"
	"context"
	"io"
	"sync"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

// Start Generation Here
const SVMKIT_SSH_STDOUT = "SVMKIT_SSH_STDOUT"
const SVMKIT_SSH_STDERR = "SVMKIT_SSH_STDERR"

func LogOutput(ctx context.Context, r io.Reader, doneCh chan<- struct{}, severity diag.Severity) {
	defer close(doneCh)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		msg := scanner.Text()
		l := p.GetLogger(ctx)
		switch severity {
		case diag.Info:
			l.InfoStatus(msg)
		case diag.Warning:
			l.WarningStatus(msg)
		case diag.Error:
			l.ErrorStatus(msg)
		default:
			l.DebugStatus(msg)
		}
	}
}

// NoopLogger satisfies the expected logger shape but doesn't actually log.
// It reads from the provided reader until EOF, discarding the output, then closes the channel.
func NoopLogger(r io.Reader, done chan struct{}) {
	defer close(done)
	_, _ = io.Copy(io.Discard, r)
}

type ConcurrentWriter struct {
	Writer io.Writer
	mu     sync.Mutex
}

func (w *ConcurrentWriter) Write(bs []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Writer.Write(bs)
}
