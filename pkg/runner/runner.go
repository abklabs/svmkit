package runner

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Command interface {
	Check() error
	Env() *EnvBuilder
	Script() string
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
		RootPath: fmt.Sprintf("/tmp/runner-%d-%d", time.Now().Unix(), rand.Int()),
	}

	p.AddString("lib.bash", LibBash)
	p.Add(PayloadFile{"run.sh", strings.NewReader(RunScript), 0755})
	p.AddReader("env", r.command.Env().Buffer())
	p.AddString("steps.sh", r.command.Script())

	d := Deployer{Payload: p, Client: r.client}

	if err := d.Deploy(); err != nil {
		return err
	}

	if err := d.Run([]string{"./run.sh"}, false, handler); err != nil {
		return err
	}

	return nil
}
