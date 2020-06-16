package node

import (
	"github.com/magiconair/properties/assert"
	"github.com/monkiato/apio-orchestrator/pkg/config"
	"testing"
)

func TestManifestFile(t *testing.T) {
	t.Run("should return manifest file using default config path", func(t *testing.T) {
		assert.Equal(t, ManifestFile("test-node"), "/var/lib/apio-orchestrator/nodes/test-node/manifest.json")
	})
}

func TestMetadataFile(t *testing.T) {
	t.Run("should return metadata file using default config path", func(t *testing.T) {
		assert.Equal(t, MetadataFile("test-node"), "/var/lib/apio-orchestrator/nodes/test-node/metadata.json")
	})
}

func TestNodeFolder(t *testing.T) {
	t.Run("should return node root folder using default config path", func(t *testing.T) {
		assert.Equal(t, NodeFolder("test-node"), "/var/lib/apio-orchestrator/nodes/test-node")
	})
}

func TestSetRootConfigPath(t *testing.T) {
	t.Run("should change default config path by custom path", func(t *testing.T) {
		SetRootConfigPath("/custom/path/")
		assert.Equal(t, ManifestFile("test-node"), "/custom/path/nodes/test-node/manifest.json")
		assert.Equal(t, MetadataFile("test-node"), "/custom/path/nodes/test-node/metadata.json")
		assert.Equal(t, NodeFolder("test-node"), "/custom/path/nodes/test-node")
		SetRootConfigPath(config.DefaultConfigPath)
	})
}
