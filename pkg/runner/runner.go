package runner

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Command interface {
	Check() error
	Env() *EnvBuilder
	AddToPayload(*Payload) error
}

func NewRunner(client *ssh.Client, cmd Command) *Runner {
	return &Runner{client: client, command: cmd}
}

type Runner struct {
	client  *ssh.Client
	command Command
}

func (r *Runner) Run(ctx context.Context, handler DeployerHandler) error {
	p := &Payload{
		RootPath:    fmt.Sprintf("/tmp/runner-%d-%d", time.Now().Unix(), rand.Int()),
		DefaultMode: 0640,
	}

	p.Add(PayloadFile{"opsh", strings.NewReader(OPSH), os.FileMode(0755)})
	p.AddString("lib.bash", LibBash)
	p.Add(PayloadFile{"run.sh", strings.NewReader(RunScript), os.FileMode(0755)})
	p.AddReader("env", r.command.Env().Buffer())

	if err := r.command.AddToPayload(p); err != nil {
		return err
	}

	d := Deployer{Payload: p, Client: r.client}

	if err := d.Deploy(); err != nil {
		return err
	}

	if err := d.Run([]string{"./run.sh"}, false, handler); err != nil {
		return err
	}

	return nil
}
