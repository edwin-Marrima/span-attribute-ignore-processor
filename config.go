package spanignoreprocessor

type Config struct {
	// IgnoredAttributes is a list of not allowed span attribute keys. Span attributes
	// that are on the list are removed.
	IgnoredAttributes AttributesConfiguration `mapstructure:"ignored_attributes"`
	//
}

type AttributesConfiguration struct {
	IncludeResources bool     `mapstructure:"include_resources"`
	Attributes       []string `mapstructure:"attributes"`
}
