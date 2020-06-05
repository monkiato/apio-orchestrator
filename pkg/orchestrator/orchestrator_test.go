package orchestrator

import (
	"github.com/monkiato/apio-orchestrator/internal/data"
	"github.com/monkiato/apio-orchestrator/pkg/persistence"
	"testing"
)

func TestNewNodeOrchestrator(t *testing.T) {
	orchestrator, err := NewNodeOrchestrator("testing", nil)

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
	conn := persistence.NewFileConnection()
	orchestrator, _ := NewNodeOrchestrator("testing1", conn)
	orchestrator.RemoveNode()

	node, err := orchestrator.CreateNode()
	t.Log(node)
	t.Log(err)

	// add collections
	orchestrator.UpdateCollections([]data.CollectionDefinition{
		{
			Name: "books",
			Fields: map[string]string{
				"title":  "string",
				"author": "string",
				"year":   "float",
			},
		},
	})

	orchestrator.AddCollection(data.CollectionDefinition{
		Name: "authors",
		Fields: map[string]string{},
	})
	orchestrator.AddCollection(data.CollectionDefinition{
		Name: "authors",
		Fields: map[string]string{},
	})
}
