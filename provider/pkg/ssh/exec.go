package ssh

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

func Exec(ctx context.Context, client *ssh.Client, command string) (string, string, error) {
	// Create a new session for running the command
	execSession, err := client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer execSession.Close()

	stdoutPipe, err := execSession.StdoutPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := execSession.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go func() {
		io.Copy(&stdoutBuf, stdoutPipe)
		close(stdoutDone)
	}()

	go func() {
		io.Copy(&stderrBuf, stderrPipe)
		close(stderrDone)
	}()

	bashCommand := fmt.Sprintf("/bin/sh -c '%s'", command)
	if err := execSession.Start(bashCommand); err != nil {
		return "", "", fmt.Errorf("failed to start command: %w", err)
	}

	if err := execSession.Wait(); err != nil {
		return "", "", fmt.Errorf("command execution failed: %w", err)
	}

	<-stdoutDone
	<-stderrDone

	return stdoutBuf.String(), stderrBuf.String(), nil
}
