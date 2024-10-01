package ssh

import (
	"context"
	"fmt"
	"net"
	"time"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	dialErrorDefault   = 1000
	dialErrorUnlimited = -1
)

// Connection represents the configuration needed to establish an SSH connection.
type Connection struct {
	User               *string  `pulumi:"user,optional"`               // The user for the SSH connection.
	Password           *string  `pulumi:"password,optional"`           // The password for the SSH connection.
	Host               *string  `pulumi:"host"`                        // The host address for the SSH connection.
	Port               *float64 `pulumi:"port,optional"`               // The port for the SSH connection. Defaults to 22.
	PrivateKey         *string  `pulumi:"privateKey,optional"`         // The private key for the SSH connection.
	PrivateKeyPassword *string  `pulumi:"privateKeyPassword,optional"` // The password for the private key if it is encrypted.
	AgentSocketPath    *string  `pulumi:"agentSocketPath,optional"`    // The SSH agent socket path.
	DialErrorLimit     *int     `pulumi:"dialErrorLimit,optional"`     // The maximum number of dial errors allowed. -1 for unlimited. Defaults to 1000.
	PerDialTimeout     *int     `pulumi:"perDialTimeout,optional"`     // The timeout for each dial attempt in seconds. Defaults to 15 seconds.
}

// Annotate adds descriptions and default values to the Connection fields.
func (c *Connection) Annotate(a infer.Annotator) {
	a.Describe(&c, "Instructions for how to connect to a remote endpoint.")
	a.Describe(&c.User, "The user that we should use for the connection.")
	a.SetDefault(&c.User, "root")
	a.Describe(&c.Password, "The password we should use for the connection.")
	a.Describe(&c.Host, "The address of the resource to connect to.")
	a.Describe(&c.Port, "The port to connect to. Defaults to 22.")
	a.SetDefault(&c.Port, 22)
	a.Describe(&c.PrivateKey, "The contents of an SSH key to use for the connection. This takes preference over the password if provided.")
	a.Describe(&c.PrivateKeyPassword, "The password to use in case the private key is encrypted.")
	a.Describe(&c.AgentSocketPath, "SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.")
	a.Describe(&c.DialErrorLimit, "Max allowed errors on trying to dial the remote host. -1 set count to unlimited. Default value is 1000.")
	a.SetDefault(&c.DialErrorLimit, dialErrorDefault)
	a.Describe(&c.PerDialTimeout, "Max number of seconds for each dial attempt. 0 implies no maximum. Default value is 15 seconds.")
	a.SetDefault(&c.PerDialTimeout, 15)
}

// Config creates and returns an SSH client configuration based on the Connection settings.
func (c *Connection) Config() (*ssh.ClientConfig, error) {
	authMethods := []ssh.AuthMethod{}

	if c.Password != nil {
		authMethods = append(authMethods, ssh.Password(*c.Password))
	}

	if c.PrivateKey != nil {
		signer, err := ssh.ParsePrivateKey([]byte(*c.PrivateKey))
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if c.AgentSocketPath != nil {
		agentConn, err := net.Dial("unix", *c.AgentSocketPath)
		if err != nil {
			return nil, err
		}
		agentClient := agent.NewClient(agentConn)
		authMethods = append(authMethods, ssh.PublicKeysCallback(agentClient.Signers))
	}

	config := &ssh.ClientConfig{
		User:            *c.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(*c.PerDialTimeout) * time.Second, // Set the timeout for the SSH connection
	}

	return config, nil
}

// Dial attempts to establish an SSH connection using the provided context and Connection settings.
func (c *Connection) Dial(ctx context.Context) (*ssh.Client, error) {
	logger := p.GetLogger(ctx)
	logger.InfoStatus("Starting SSH dial process")

	config, err := c.Config()
	if err != nil {
		logger.ErrorStatus(fmt.Sprintf("failed to get SSH config: %v", err))
		return nil, fmt.Errorf("failed to get SSH config: %w", err)
	}

	var client *ssh.Client
	var dialErr error

	for i := 0; i < *c.DialErrorLimit || *c.DialErrorLimit == dialErrorUnlimited; i++ {
		client, dialErr = ssh.Dial("tcp", fmt.Sprintf("%s:%d", *c.Host, int(*c.Port)), config)
		if dialErr == nil {
			logger.InfoStatus("SSH connection established successfully")
			return client, nil
		}

		logger.WarningStatus(fmt.Sprintf("failed to dial SSH (attempt %d): %v", i+1, dialErr))

		select {
		case <-ctx.Done():
			logger.ErrorStatus(fmt.Sprintf("context deadline exceeded: %v", ctx.Err()))
			return nil, fmt.Errorf("context deadline exceeded: %w", ctx.Err())
		default:
		}

		if ctx.Err() == context.DeadlineExceeded {
			logger.WarningStatus("dial context deadline exceeded, retrying...")
			continue
		}
	}

	logger.ErrorStatus(fmt.Sprintf("failed to establish SSH connection after %d attempts: %v", *c.DialErrorLimit, dialErr))
	return nil, fmt.Errorf("failed to establish SSH connection after %d attempts: %w", *c.DialErrorLimit, dialErr)
}
