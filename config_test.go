package spanignoreprocessor

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestLoadingConfig(t *testing.T) {
	testCase := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			id: component.NewIDWithName(typeStr, ""),
			expected: &Config{
				IgnoredAttributes: []string{"token"},
			},
		},
	}
	for _, tt := range testCase {
		cm, err := confmaptest.LoadConf(filepath.Join("test_artifacts", "config1.yaml"))
		require.NoError(t, err)

		factory := NewFactory()
		cfg := factory.CreateDefaultConfig()

		sub, err := cm.Sub(tt.id.String())
		require.NoError(t, err)
		require.NoError(t, component.UnmarshalConfig(sub, cfg))

		assert.NoError(t, component.ValidateConfig(cfg))
		assert.Equal(t, tt.expected, cfg)

	}
}
