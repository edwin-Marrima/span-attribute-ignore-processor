package spanignoreprocessor

type Config struct {
	// IgnoredAttributes is a list of not allowed span attribute keys. Span attributes
	// that are on the list are removed.
	IgnoredAttributes []string `mapstructure:"ignored_attributes"`
	//
}
