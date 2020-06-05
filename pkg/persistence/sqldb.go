package persistence

import (
	"github.com/jinzhu/gorm"
	"github.com/monkiato/apio-orchestrator/internal/models"
)

type SqlConnection struct {
	db *gorm.DB
}

//NewSqlConnection create a new SqlConnection instance with a reference to the persistence passed by param
func NewSqlConnection(db *gorm.DB) Connection {
	return &SqlConnection{
		db: db,
	}
}

//GetNode looks for a specific node data in the storage/DB
func (conn *SqlConnection) GetNode(nodeId string) (*models.Node, bool) {
	var node models.Node
	if conn.db.Where("node_id = ?", nodeId).First(&node).RecordNotFound() {
		return nil, false
	}
	return &node, true
}

//NodeExists returns if the node exists in the DB
func (conn *SqlConnection) NodeExists(nodeId string) bool {
	_, exists := conn.GetNode(nodeId)
	return exists
}

//CreateNode creates new node in the persistence
func (conn *SqlConnection) CreateNode(node *models.Node) error {
	return conn.db.Create(node).Error
}


//UpdateNode updates existing node in the DB
func (conn *SqlConnection) UpdateNode(node *models.Node) error {
	return conn.db.Save(&node).Error
}

//RemoveNode removes node from DB
func (conn *SqlConnection) RemoveNode(nodeId string) error {
	return conn.db.Delete(&models.Node{
		NodeId: nodeId,
	}).Error
}
