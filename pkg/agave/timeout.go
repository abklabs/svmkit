package agave

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type TimeoutConfig struct {
	RpcServiceTimeout *int `pulumi:"rpcServiceTimeout,optional"`
}

func (t *TimeoutConfig) Env() *runner.EnvBuilder {
	e := runner.NewEnvBuilder()

	e.SetIntP("RPC_SERVICE_TIMEOUT", t.RpcServiceTimeout)

	return e
}
