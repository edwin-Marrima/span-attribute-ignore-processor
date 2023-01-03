package spanignoreprocessor

import (
	"context"
	"fmt"
	"regexp"

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
	//list of not allowed span events
	ignoredEvents map[string]*regexp.Regexp
	config        *Config
	// Logger
	logger *zap.Logger
	// Next trace consumer in line
	next consumer.Traces
}

func NewProofreader(ctx context.Context, config *Config, logger *zap.Logger, next consumer.Traces) (*proofreader, error) {
	compiledRegex, err := buildRegex(config.IgnoredEvents)
	if err != nil {
		return nil, fmt.Errorf("error creating `proofreader` processor: %w", err)
	}
	return &proofreader{
		config:           config,
		logger:           logger,
		next:             next,
		ignoreAttributes: buildIgnoreAttributes(config.IgnoredAttributes.Attributes),
		ignoredEvents:    compiledRegex,
	}, nil
}

func buildIgnoreAttributes(elements []string) map[string]struct{} {
	ignoreAttributes := map[string]struct{}{}
	for _, v := range elements {
		ignoreAttributes[v] = struct{}{}
	}
	return ignoreAttributes
}
func buildRegex(expressions []string) (map[string]*regexp.Regexp, error) {
	regex := make(map[string]*regexp.Regexp)
	for _, expression := range expressions {
		re, err := regexp.Compile(expression)
		if err != nil {
			return nil, fmt.Errorf("Error compiling regex in block list: %s", err.Error())
		}
		regex[expression] = re
	}
	return regex, nil
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
		//When IncludeResources is false resources attributes are ignored
		if proofr.config.IgnoredAttributes.IncludeResources {
			resourceAttributes := rs.Resource().Attributes()
			processAttributes(proofr.ignoreAttributes, resourceAttributes)
		}
		for j := 0; j < rs.ScopeSpans().Len(); j++ {
			ils := rs.ScopeSpans().At(j)
			for k := 0; k < ils.Spans().Len(); k++ {
				span := ils.Spans().At(k)
				spanAttrs := span.Attributes()
				processAttributes(proofr.ignoreAttributes, spanAttrs)
				spanEvent := span.Events()
				//iterate over the events in order to eliminate the events that
				//satisfy(match) the previously provided regular expressions
				spanEvent.RemoveIf(func(se ptrace.SpanEvent) bool {
					for _, re := range proofr.ignoredEvents {
						// verify if regex match event name
						if re.MatchString(se.Name()) {
							return true
						}
					}
					return false
				})
			}
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
