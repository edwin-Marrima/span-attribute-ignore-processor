package spanignoreprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestFactoryType(t *testing.T) {
	factory := NewFactory()
	assert.Equal(t, factory.Type(), component.Type(typeStr))
}

func TestDefaultConfiguration(t *testing.T) {
	factory := NewFactory()
	defaultConfiguration := factory.CreateDefaultConfig()
	expectedDefaultConfiguration := &Config{
		IgnoredAttributes: AttributesConfiguration{
			IncludeResources: true,
		},
	}
	assert.Equal(t, expectedDefaultConfiguration, defaultConfiguration)
}

func TestCreateTracesProcessor(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)

	tracesProcessor, err := factory.CreateTracesProcessor(context.Background(), processortest.NewNopCreateSettings(), oCfg, consumertest.NewNop())
	require.Nil(t, err)
	assert.NotNil(t, tracesProcessor)
}
