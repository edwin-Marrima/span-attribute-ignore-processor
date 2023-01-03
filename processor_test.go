package spanignoreprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processortest"
)

type testCase struct {
	name               string
	serviceName        string
	inputAttributes    map[string]interface{}
	expectedAttributes map[string]interface{}
}

func runIndividualTestCase(t *testing.T, tt testCase, tp processor.Traces) {
	t.Run(tt.name, func(t *testing.T) {
		td := generateTraceData(tt.name, tt.inputAttributes)
		assert.NoError(t, tp.ConsumeTraces(context.Background(), td))
		// Ensure that the modified `td` has the attributes sorted:
		sortAttributes(td)
		require.Equal(t, generateTraceData(tt.name, tt.expectedAttributes), td)
	})
}
func generateTraceData(spanName string, spanAttributes map[string]interface{}) ptrace.Traces {
	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()

	scopeSpan := rs.ScopeSpans().AppendEmpty()
	span := scopeSpan.Spans().AppendEmpty()
	span.SetName(spanName)
	span.SetTraceID([16]byte{1, 2, 3, 4})

	span.Attributes().FromRaw(spanAttributes)

	span.Attributes().Sort()

	return td
}
func sortAttributes(td ptrace.Traces) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		rs.Resource().Attributes().Sort()
		ilss := rs.ScopeSpans()
		for j := 0; j < ilss.Len(); j++ {
			spans := ilss.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				spans.At(k).Attributes().Sort()
			}
		}
	}
}

func TestIgnoreSpans(t *testing.T) {
	testCases := []struct {
		test   testCase
		config *Config
	}{
		{
			test: testCase{
				name:        "Remove span whose Key is in IgnoredAttributes property",
				serviceName: "admin_service",
				inputAttributes: map[string]interface{}{
					"account.id":       "007",
					"http.status_code": 200,
					"account.password": "AKdhcjs^&xva",
				},
				expectedAttributes: map[string]interface{}{
					"account.id":       "007",
					"http.status_code": 200,
				},
			},
			config: &Config{
				IgnoredAttributes: []string{"account.password"},
			},
		},
	}
	for _, v := range testCases {
		factory := NewFactory()
		tp, err := factory.CreateTracesProcessor(context.Background(), processortest.NewNopCreateSettings(), v.config, consumertest.NewNop())
		require.Nil(t, err)
		require.NotNil(t, tp)
		runIndividualTestCase(t, v.test, tp)
	}
}
