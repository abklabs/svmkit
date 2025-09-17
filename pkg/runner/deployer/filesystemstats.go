package deployer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
)

type FileSystemStats struct {
	Path       string `json:"path"`
	BlockSize  uint64 `json:"blockSize"`
	FreeBlocks uint64 `json:"freeBlocks"`
}

func (fss *FileSystemStats) FreeBytes() uint64 {
	return fss.BlockSize * fss.FreeBlocks
}

func GetFileSystemStats(client *ssh.Client, path string) (stats FileSystemStats, err error) {
	if client == nil {
		return FileSystemStats{}, fmt.Errorf("SSH client cannot be nil")
	}

	execSession, err := client.NewSession()
	if err != nil {
		return
	}

	defer func() {
		if closeErr := execSession.Close(); closeErr != io.EOF {
			err = errors.Join(err, closeErr)
		}
	}()

	cmd := fmt.Sprintf("stat -f -c '{\"path\": %q, \"blockSize\": %%S, \"freeBlocks\": %%a}' %q", path, path)

	out, err := execSession.CombinedOutput(cmd)
	if err != nil {
		return stats, fmt.Errorf("remote stat failed (output: %q): %w", out, err)
	}

	if err := json.Unmarshal(bytes.TrimSpace(out), &stats); err != nil {
		return stats, fmt.Errorf("invalid JSON from remote stat (got %v): %w", string(out), err)
	}
	return
}
