package node

import "monkiato/apio-orchestrator/internal/data"

//Metadata is the info related to a single Apio node
type Metadata struct {
	Name        string                      `json:"name"`
	collections []data.CollectionDefinition `json:"collections"`
	Active      bool                        `json:"active"`
}
