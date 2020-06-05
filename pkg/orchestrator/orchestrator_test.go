package orchestrator

import (
	"monkiato/apio-orchestrator/internal/data"
	"monkiato/apio-orchestrator/pkg/node"
	"testing"
)

func TestNewNodeOrchestrator(t *testing.T) {
	orchestrator, err := NewNodeOrchestrator(&node.Metadata{
		Name:   "testing",
		Active: true,
	})

	if err != nil {
		t.Fatalf("unable to create NodeOrchestrator instance")
	}
	if orchestrator.node == nil {
		t.Fatalf("invalid nil node")
	}
	if orchestrator.dockerClient == nil {
		t.Fatalf("docker client not initialized")
	}
}

func TestNodeOrchestrator_CreateNode(t *testing.T) {
	orchestrator, _ := NewNodeOrchestrator(&node.Metadata{
		Name: "testing",
		Collections: []data.CollectionDefinition{
			{
				Name: "books",
				Fields: map[string]string{
					"title": "string",
					"author": "string",
					"year": "float",
				},
			},
		},
		Active: true,
	})
	err := orchestrator.CreateNode()
	t.Log(err)
}
