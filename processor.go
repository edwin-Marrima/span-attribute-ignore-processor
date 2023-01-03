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
		ignoreAttributes: buildIgnoreAttributes(config.IgnoredAttributes.Attributes),
	}
}

func buildIgnoreAttributes(elements []string) map[string]struct{} {
	ignoreAttributes := map[string]struct{}{}
	for _, v := range elements {
		ignoreAttributes[v] = struct{}{}
	}
	return ignoreAttributes
}

func (proofr *proofreader) processTraces(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {

	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		rs := traces.ResourceSpans().At(i)

		// span-id
		// ilss := rs.ScopeSpans()
		// ils := ilss.At(0)
		// spans := ils.Spans()
		// span := spans.At(0)
		// span.SpanID()
		// span.SetSpanID()
		// span.SetTraceID()
		//

		for j := 0; j < rs.ScopeSpans().Len(); j++ {
			ils := rs.ScopeSpans().At(j)
			for k := 0; k < ils.Spans().Len(); k++ {
				span := ils.Spans().At(k)
				spanAttrs := span.Attributes()
				processAttributes(proofr.ignoreAttributes, spanAttrs)
			}
		}
		//When IncludeResources is false resources attributes are ignored
		if proofr.config.IgnoredAttributes.IncludeResources {
			resourceAttributes := rs.Resource().Attributes()
			processAttributes(proofr.ignoreAttributes, resourceAttributes)
		}

	}
	return traces, nil
}

func processAttributes(ignoreAttributes map[string]struct{}, attributes pcommon.Map) {
	//attributes is passed by reference (pcommon.Map is a pointer), reason why this method doesn't return any value
	attributes.RemoveIf(func(key string, v pcommon.Value) bool {
		//attributes are removed only when the list of attributes provided in config.yaml exists in list(slice) of span or resource attributes
		if _, ok := ignoreAttributes[key]; ok {
			return true
		}
		return false
	})
}

// Processors which modify the input data MUST set this flag to true
func (proofr *proofreader) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}
