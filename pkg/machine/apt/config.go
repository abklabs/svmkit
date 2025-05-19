package apt

type Config struct {
	ExcludeDefaultSources *bool `pulumi:"excludeDefaultSources,optional"`
}
