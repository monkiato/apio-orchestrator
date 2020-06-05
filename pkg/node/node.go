package node

import (
	"fmt"
	"github.com/monkiato/apio-orchestrator/internal/data"
	"github.com/monkiato/apio-orchestrator/pkg/config"
	"os"
)

//Metadata is the info related to a single Apio node
type Metadata struct {
	Name        string                      `json:"name"`
	Collections []data.CollectionDefinition `json:"collections"`
	Active      bool                        `json:"active"`
}

//ManifestFile indicates the absolute path where the manifest file is located for the specified node
func ManifestFile(nodeId string) string {
	return fmt.Sprintf("%s/manifest.json", NodeFolder(nodeId))
}

//MetadataFile indicates the absolute path where the metadata file is located for the specified node
func MetadataFile(nodeId string) string {
	return fmt.Sprintf("%s/metadata.json", NodeFolder(nodeId))
}

//NodeFolder indicates the absolute folder used for node information
func NodeFolder(nodeId string) string {
	return fmt.Sprintf("%s%s", config.RootNodeFolder, nodeId)
}

func CreateNodeFolder(nodeId string) error {
	return os.MkdirAll(NodeFolder(nodeId), 0777)
}