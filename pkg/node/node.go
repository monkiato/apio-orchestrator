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

//DockerConfig specify details required to create and manage the container related to the node
type DockerConfig struct {
	NetworkName         string
	DomainName          string
	ContainerNamePrefix string
	MongoDbHost         string
	MongoDbName         string
}

var rootConfigPath string = config.DefaultConfigPath

//SetRootConfigPath the setter must be used to modify the root path if not using the default
func SetRootConfigPath(configPath string) {
	rootConfigPath = configPath
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
	return fmt.Sprintf("%s%s", rootConfigPath+config.NodeFolder, nodeId)
}

//CreateNodeConfigFolder ensure the folder structure is created for the specified node
func CreateNodeConfigFolder(nodeId string) error {
	return os.MkdirAll(NodeFolder(nodeId), 0777)
}
