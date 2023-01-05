package spanignoreprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "proofreader"
)

func createDefaultConfig() component.Config {
	return &Config{
		IgnoredAttributes: AttributesConfiguration{
			IncludeResources: true,
		},
	}
}

func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, component.StabilityLevelDevelopment),
	)
}

func createTracesProcessor(
	ctx context.Context,
	set processor.CreateSettings,
	cfg component.Config,
	next consumer.Traces,
) (processor.Traces, error) {

	oCfg := cfg.(*Config)

	proofreader, err := NewProofreader(ctx, oCfg, set.Logger, next)
	if err != nil {
		return nil, err
	}
	return processorhelper.NewTracesProcessor(
		ctx,
		set,
		cfg,
		next,
		proofreader.processTraces,
		processorhelper.WithCapabilities(proofreader.Capabilities()),
	)
}
