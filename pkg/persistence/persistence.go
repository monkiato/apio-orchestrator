package persistence

import (
	"github.com/monkiato/apio-orchestrator/internal/models"
)

type Connection interface {
	//GetNode looks for a specific node data in the storage/DB
	GetNode(nodeId string) (*models.Node, bool)
	//NodeExists returns if the node exists in the DB
	NodeExists(nodeId string) bool
	//CreateNode creates new node in the persistence
	CreateNode(node *models.Node) error
	//UpdateNode updates existing node in the DB
	UpdateNode(node *models.Node) error
	//RemoveNode removes node from DB
	RemoveNode(nodeId string) error
}
