package persistence

import (
	"encoding/json"
	"fmt"
	"github.com/monkiato/apio-orchestrator/internal/models"
	"github.com/monkiato/apio-orchestrator/pkg/node"
	"io/ioutil"
	"os"
	"time"
)

type FileConnection struct {
	configPath string
}

//NewFileConnection create a new FileConnection instance
func NewFileConnection(configPath string) Connection {
	return &FileConnection{
		configPath: configPath,
	}
}

//GetNode looks for a specific node data in the storage
func (conn *FileConnection) GetNode(nodeId string) (*models.Node, bool) {
	data, error := ioutil.ReadFile(node.MetadataFile(nodeId))
	if error != nil {
		//TODO: log error
		return nil, false
	}
	var node models.Node
	if json.Unmarshal(data, &node) != nil {
		//TODO: log error
		return nil, false
	}
	return &node, true
}

//NodeExists returns if the node exists in the persisted data
func (conn *FileConnection) NodeExists(nodeId string) bool {
	info, err := os.Stat(node.MetadataFile(nodeId))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//CreateNode creates new node in the persistence folder
func (conn *FileConnection) CreateNode(newNode *models.Node) error {
	if conn.NodeExists(newNode.NodeId) {
		return fmt.Errorf("node with ID '%s' already exists", newNode.NodeId)
	}
	newNode.CreatedAt = time.Now()
	return conn.UpdateNode(newNode)
}

//UpdateNode updates existing node in the persisted data folder
func (conn *FileConnection) UpdateNode(_node *models.Node) error {
	_node.UpdatedAt = time.Now()
	data, _ := json.MarshalIndent(_node, "", "    ")
	return ioutil.WriteFile(node.MetadataFile(_node.NodeId), data, 0644)
}

//RemoveNode removes node from the persisted data folder
func (conn *FileConnection) RemoveNode(nodeId string) error {
	return os.Remove(node.MetadataFile(nodeId))
}
