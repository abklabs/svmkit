package ssh

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Upload(ctx context.Context, client *ssh.Client, content []byte, path string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	// Create the directory structure if it doesn't exist
	dir := filepath.Dir(path)
	if err := sftpClient.MkdirAll(dir); err != nil {
		return fmt.Errorf("failed to create remote directory %s: %w", dir, err)
	}

	// Create the remote file
	remoteFile, err := sftpClient.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create remote file %s: %w", path, err)
	}
	defer remoteFile.Close()

	// Write the content to the remote file
	if _, err := remoteFile.Write(content); err != nil {
		return fmt.Errorf("failed to write to remote file %s: %w", path, err)
	}

	return nil
}
