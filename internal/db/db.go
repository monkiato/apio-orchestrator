package db

import "monkiato/apio-orchestrator/internal/node"

//GetNode looks for a specific node data in the storage/DB
func GetNode(nodeId string) (*node.Metadata, bool) {
	//TODO: look for nodeID in the DB
	return nil, false
}
