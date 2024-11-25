package agave

type StartupPolicy struct {
	WaitForRPCHealth *bool `pulumi:"waitForRPCHealth,optional"`
}
