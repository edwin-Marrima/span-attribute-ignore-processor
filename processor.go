package spanignoreprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type proofreader struct {
	_ struct{}

	//since processors do a lot of contains checks, iterating over the list often has a serious impact on performance
	// it explains the usage of map[string]struct{} type rather than slice.
	//Using an empty struct{} here has the advantage that it doesn't require any additional space and Go's internal map
	//type is optimized for that kind of values.
	//list of not allowed span attribute keys
	ignoreAttributes map[string]struct{}
	config           *Config
	// Logger
	logger *zap.Logger
	// Next trace consumer in line
	next consumer.Traces
}

func NewProofreader(ctx context.Context, config *Config, logger *zap.Logger, next consumer.Traces) *proofreader {
	return &proofreader{
		config:           config,
		logger:           logger,
		next:             next,
		ignoreAttributes: buildIgnoreAttributes(config.IgnoredAttributes),
	}
}

func buildIgnoreAttributes(elements []string) map[string]struct{} {
	ignoreAttributes := map[string]struct{}{}
	for _, v := range elements {
		ignoreAttributes[v] = struct{}{}
	}
	return ignoreAttributes
}

func (ig *proofreader) processTraces(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		rs := traces.ResourceSpans().At(i)

		resourceAttributes := rs.Resource().Attributes()
		resourceAttributes.RemoveIf(func(key string, v pcommon.Value) bool {
			if _, ok := ig.ignoreAttributes[key]; ok {
				return true
			}
			return false
		})
	}

	return traces, nil
}
