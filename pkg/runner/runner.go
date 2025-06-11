package runner

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/abklabs/svmkit/pkg/runner/deployer"

	"golang.org/x/crypto/ssh"
)

type Command interface {
	Check() error
	Env() *EnvBuilder
	AddToPayload(*Payload) error
	Config() *Config
}

func NewRunner(client *ssh.Client, cmd Command) *Runner {
	return &Runner{client: client, command: cmd}
}

type Runner struct {
	client  *ssh.Client
	command Command
}

func PrepareCommandPayload(p *Payload, command Command) error {
	p.Add(PayloadFile{Path: "opsh", Reader: strings.NewReader(OPSH), Mode: 0755})
	p.AddString("lib.bash", LibBash)
	p.Add(PayloadFile{Path: "run.sh", Reader: strings.NewReader(RunScript), Mode: 0755})
	p.AddReader("env", command.Env().Buffer())

	if err := command.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

func (r *Runner) Run(ctx context.Context, handler deployer.DeployerHandler, statusCallback deployer.ProgressStatusCallback) error {
	p := &Payload{
		RootPath:    fmt.Sprintf("/tmp/runner-%d-%d", time.Now().Unix(), rand.Int()),
		DefaultMode: 0640,
	}

	if err := PrepareCommandPayload(p, r.command); err != nil {
		return err
	}

	keepPayload := false

	if c := r.command.Config(); c != nil {
		if c.KeepPayload != nil {
			keepPayload = *c.KeepPayload
		}
	}

	d := deployer.SSH{Payload: p, Client: r.client, KeepPayload: keepPayload}
	if err := d.Deploy(statusCallback); err != nil {
		return err
	}

	if err := d.Run([]string{"./run.sh"}, handler); err != nil {
		return err
	}

	return nil
}
