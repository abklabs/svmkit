package ssh

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func Ready(ctx context.Context, client *ssh.Client) error {
	var backoff = time.Second

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context deadline exceeded: %w", ctx.Err())
		default:
			err := Exec(ctx, client, "echo 'ping'")
			if err == nil {
				return nil
			}
			time.Sleep(backoff)
			backoff *= 2
			if backoff > time.Minute {
				backoff = time.Minute
			}
		}
	}
}
