package proxy

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetTraefikLabels(t *testing.T) {
	t.Run("should return traefik label", func(t *testing.T) {
		labels := GetTraefikLabels(
			"test-node",
			"localhost",
			"test-network",
			8080)
		assert.Equal(t, len(labels), 4, "unexpected labels length")
		assert.Equal(t, keyExists(labels, "traefik.http.routers.test-node.rule"), true, "http router key not found in map")
		assert.Equal(t, keyExists(labels, "traefik.docker.network"), true, "docker network key not found in map")
		assert.Equal(t, keyExists(labels, "traefik.enable"), true, "enable key not found in map")
		assert.Equal(t, keyExists(labels, "traefik.http.services.test-node.loadbalancer.server.port"), true, "loadbalancer port key not found in map")

		assert.Equal(t, labels["traefik.http.routers.test-node.rule"], "Host(`localhost`)", "non-matching value")
		assert.Equal(t, labels["traefik.docker.network"], "test-network", "non-matching value")
		assert.Equal(t, labels["traefik.enable"], "true", "non-matching value")
		assert.Equal(t, labels["traefik.http.services.test-node.loadbalancer.server.port"], "8080", "non-matching value")
	})
}

func keyExists(m map[string]string, key string) bool {
	_, exists := m[key]
	return exists
}
