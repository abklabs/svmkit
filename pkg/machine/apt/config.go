package apt

type Config struct {
	Sources               *Sources `pulumi:"sources,optional"`
	ExcludeDefaultSources *bool    `pulumi:"excludeDefaultSources,optional"`
}
