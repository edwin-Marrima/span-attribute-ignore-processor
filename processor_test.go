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
	name                       string
	serviceName                string
	spanInputAttributes        map[string]interface{}
	spanExpectedAttributes     map[string]interface{}
	resourceInputAttributes    map[string]interface{}
	resourceExpectedAttributes map[string]interface{}
}

func runIndividualTestCase(t *testing.T, tt testCase, tp processor.Traces) {
	t.Run(tt.name, func(t *testing.T) {
		td := generateTraceData(tt.name, tt.spanInputAttributes, tt.resourceInputAttributes)
		assert.NoError(t, tp.ConsumeTraces(context.Background(), td))
		// Ensure that the modified `td` has the attributes sorted:
		sortAttributes(td)
		require.Equal(t, generateTraceData(tt.name, tt.spanExpectedAttributes, tt.resourceExpectedAttributes), td)
	})
}
func generateTraceData(spanName string, spanAttributes map[string]interface{}, resourceAttributes map[string]interface{}) ptrace.Traces {
	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	if len(resourceAttributes) > 0 {
		rs.Resource().Attributes().FromRaw(resourceAttributes)
		rs.Resource().Attributes().Sort()
	}
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
				name:        "Remove span attribute whose Key is listed in IgnoredAttributes property",
				serviceName: "admin_service",
				spanInputAttributes: map[string]interface{}{
					"account.id":       "007",
					"http.status_code": 200,
					"account.password": "AKdhcjs^&xva",
				},
				spanExpectedAttributes: map[string]interface{}{
					"account.id":       "007",
					"http.status_code": 200,
				},
			},
			config: &Config{
				IgnoredAttributes: AttributesConfiguration{
					IncludeResources: false,
					Attributes:       []string{"account.password"},
				},
			},
		},
		{
			test: testCase{
				name:        "Remove resource attribute whose Key is listed in IgnoredAttributes property",
				serviceName: "admin_service",
				resourceInputAttributes: map[string]interface{}{
					"service.namespace":   "mpt-001",
					"service.instance.id": "0x3467Sdfjk",
					"service.owner.name":  "alpha-team",
					"service.owner.hash":  "GfghjW$dshjkl32UYK",
				},
				resourceExpectedAttributes: map[string]interface{}{
					"service.namespace":   "mpt-001",
					"service.instance.id": "0x3467Sdfjk",
					"service.owner.name":  "alpha-team",
				},
			},
			config: &Config{
				IgnoredAttributes: AttributesConfiguration{
					IncludeResources: true,
					Attributes:       []string{"service.owner.hash"},
				},
			},
		},
		{
			test: testCase{
				name:        "Ignore resource attribute when IncludeResources is false",
				serviceName: "admin_service",
				resourceInputAttributes: map[string]interface{}{
					"service.namespace":   "mpt-001",
					"service.instance.id": "0x3467Sdfjk",
					"service.owner.name":  "alpha-team",
					"service.owner.hash":  "GfghjW$dshjkl32UYK",
				},
				resourceExpectedAttributes: map[string]interface{}{
					"service.namespace":   "mpt-001",
					"service.instance.id": "0x3467Sdfjk",
					"service.owner.name":  "alpha-team",
					"service.owner.hash":  "GfghjW$dshjkl32UYK",
				},
			},
			config: &Config{
				IgnoredAttributes: AttributesConfiguration{
					IncludeResources: false,
					Attributes:       []string{"service.owner.hash"},
				},
			},
		},
		{
			test: testCase{
				name:        "Remove span and resource attributes",
				serviceName: "admin_service",
				resourceInputAttributes: map[string]interface{}{
					"service.owner.name": "alpha-team",
					"service.owner.hash": "GfghjW$dshjkl32UYK",
				},
				resourceExpectedAttributes: map[string]interface{}{
					"service.owner.name": "alpha-team",
				},
				spanInputAttributes: map[string]interface{}{
					"account.id":       "007",
					"http.status_code": 200,
					"account.password": "AKdhcjs^&xva",
				},
				spanExpectedAttributes: map[string]interface{}{
					"account.id":       "007",
					"http.status_code": 200,
				},
			},
			config: &Config{
				IgnoredAttributes: AttributesConfiguration{
					IncludeResources: true,
					Attributes: []string{
						"service.owner.hash",
						"account.password",
					},
				},
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
