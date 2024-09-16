package ssh

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"golang.org/x/crypto/ssh"
)

func Exec(ctx context.Context, client *ssh.Client, command string) error {
	// Create a new session for running the command
	execSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer execSession.Close()

	stdoutPipe, err := execSession.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := execSession.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go LogOutput(ctx, stdoutPipe, stdoutDone, diag.Info)
	go LogOutput(ctx, stderrPipe, stderrDone, diag.Error)

	bashCommand := fmt.Sprintf("/bin/sh -c '%s'", command)
	if err := execSession.Start(bashCommand); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	if err := execSession.Wait(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	<-stdoutDone
	<-stderrDone

	return nil
}
