package spanignoreprocessor

type Config struct {
	IgnoredAttributes AttributesConfiguration `mapstructure:"ignored_attributes"`
	//blocked_events is a list of regular expressions for blocking events
	//whose names match the provided regular expression
	IgnoredEvents []string `mapstructure:"ignored_events"`
}

type AttributesConfiguration struct {
	//IncludeResources is a boolean value that determines whether the processor will remove
	//resources Attributes listed in Attributes property
	IncludeResources bool `mapstructure:"include_resources"`
	// IgnoredAttributes is a list of not allowed span attribute keys. Span attributes
	// that are on the list are removed.
	Attributes []string `mapstructure:"attributes"`
}
