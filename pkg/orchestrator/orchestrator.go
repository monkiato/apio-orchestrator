package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/monkiato/apio-orchestrator/internal/data"
	"github.com/monkiato/apio-orchestrator/internal/models"
	"github.com/monkiato/apio-orchestrator/internal/tools/proxy"
	"github.com/monkiato/apio-orchestrator/pkg/config"
	"github.com/monkiato/apio-orchestrator/pkg/node"
	"github.com/monkiato/apio-orchestrator/pkg/persistence"
	"io/ioutil"
	"os"
	"time"
)

//NodeOrchestrator provides functionality to interact with a single apio node (create, remove, edit)
type NodeOrchestrator struct {
	node             *node.Metadata
	nodeDockerConfig *node.DockerConfig
	dockerClient     *client.Client
	persistence      persistence.Connection
}

//NewNodeOrchestrator create new NodeOrchestrator instance to manipulate the specified nodeId
func NewNodeOrchestrator(nodeId string, nodeDockerConfig *node.DockerConfig, persistence persistence.Connection) (*NodeOrchestrator, error) {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	var collections []data.CollectionDefinition
	data, err := ioutil.ReadFile(node.ManifestFile(nodeId))
	if err == nil {
		if err := json.Unmarshal(data, &collections); err != nil {
			return nil, err
		}
	}

	orchestrator := &NodeOrchestrator{
		node: &node.Metadata{
			Name:        nodeId,
			Collections: collections,
			Active:      true,
		},
		nodeDockerConfig: nodeDockerConfig,
		dockerClient:     dockerClient,
		persistence:      persistence,
	}
	if err := orchestrator.fetchImage(); err != nil {
		return nil, err
	}
	return orchestrator, nil
}

//CreateNode create new Apio node in docker. Node metadata is also updated.
func (orchestrator *NodeOrchestrator) CreateNode() (*models.Node, error) {
	if err := node.CreateNodeConfigFolder(orchestrator.node.Name); err != nil {
		return nil, err
	}
	if err := orchestrator.createNodeManifest(); err != nil {
		return nil, err
	}
	if err := orchestrator.ensureNetworkIsCreated(); err != nil {
		return nil, err
	}
	containerId, err := orchestrator.createContainer()
	if err != nil {
		return nil, err
	}

	node := &models.Node{
		NodeId:       orchestrator.node.Name,
		ContainerId:  containerId,
		NodeFolder:   node.NodeFolder(orchestrator.node.Name),
		NodeManifest: node.ManifestFile(orchestrator.node.Name),
	}

	// save container in persistence
	if err := orchestrator.persistence.CreateNode(node); err != nil {
		// remove container created previously
		orchestrator.removeContainer(containerId)
		return nil, err
	}
	return node, nil
}

//StartNode will start the docker container associated to the Apio node.
//Node metadata is also updated.
func (orchestrator *NodeOrchestrator) StartNode() error {
	nodeData, exists := orchestrator.persistence.GetNode(orchestrator.node.Name)
	if !exists {
		return fmt.Errorf("node metadata not found")
	}
	if err := orchestrator.startContainer(nodeData.ContainerId); err != nil {
		return err
	}
	return nil
}

//StopNode will stop the docker container associated to the Apio node.
//Node metadata is also updated.
func (orchestrator *NodeOrchestrator) StopNode() error {
	nodeData, exists := orchestrator.persistence.GetNode(orchestrator.node.Name)
	if !exists {
		return fmt.Errorf("node metadata not found")
	}
	if err := orchestrator.stopContainer(nodeData.ContainerId); err != nil {
		return err
	}
	return nil
}

//RemoveNode will remove the docker container associated to the Apio node.
//Node metadata is also updated.
func (orchestrator *NodeOrchestrator) RemoveNode() error {
	nodeData, exists := orchestrator.persistence.GetNode(orchestrator.node.Name)
	if !exists {
		return fmt.Errorf("node metadata not found")
	}
	//remove from docker
	if err := orchestrator.dockerClient.ContainerRemove(
		context.Background(),
		nodeData.ContainerId,
		types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
		return err
	}

	//remove metadata
	if err := orchestrator.persistence.RemoveNode(orchestrator.node.Name); err != nil {
		return err
	}

	//remove node folder
	return os.RemoveAll(node.NodeFolder(orchestrator.node.Name))
}

//UpdateCollections override existing node collections with the new collections passed by param.
//Node metadata is also updated.
func (orchestrator *NodeOrchestrator) UpdateCollections(collections []data.CollectionDefinition) error {
	nodeData, exists := orchestrator.persistence.GetNode(orchestrator.node.Name)
	if !exists {
		return fmt.Errorf("node not found")
	}

	//update manifest
	orchestrator.node.Collections = collections
	if err := orchestrator.updateNodeManifest(); err != nil {
		return err
	}

	//restart container
	return orchestrator.restartContainer(nodeData.ContainerId)
}

//AddCollection add a new collection to the node, existing collection are not modified and remain associated to the node.
//Node metadata is also updated.
func (orchestrator *NodeOrchestrator) AddCollection(newCollection data.CollectionDefinition) error {
	nodeData, exists := orchestrator.persistence.GetNode(orchestrator.node.Name)
	if !exists {
		return fmt.Errorf("node not found")
	}

	for _, collection := range orchestrator.node.Collections {
		if collection.Name == newCollection.Name {
			return fmt.Errorf("collection with name '%s' already exists", newCollection.Name)
		}
	}

	//update manifest
	orchestrator.node.Collections = append(orchestrator.node.Collections, newCollection)
	if err := orchestrator.updateNodeManifest(); err != nil {
		return err
	}

	//restart container
	return orchestrator.restartContainer(nodeData.ContainerId)
}

//RemoveCollection remove an existing collection in the node.
//Node metadata is also updated.
func (orchestrator *NodeOrchestrator) RemoveCollection(collectionName string) error {
	nodeData, exists := orchestrator.persistence.GetNode(orchestrator.node.Name)
	if !exists {
		return fmt.Errorf("node not found")
	}

	var found bool
	for index, collection := range orchestrator.node.Collections {
		if collection.Name == collectionName {
			orchestrator.node.Collections = append(orchestrator.node.Collections[:index], orchestrator.node.Collections[index+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	//update manifest
	if err := orchestrator.updateNodeManifest(); err != nil {
		return err
	}

	//restart container
	return orchestrator.restartContainer(nodeData.ContainerId)
}

//AddField add a single field to the specified collection, field name and type must be specified.
//Node metadata is also updated.
func (orchestrator *NodeOrchestrator) AddField(collectionName string, field string, fieldType string) error {
	nodeData, exists := orchestrator.persistence.GetNode(orchestrator.node.Name)
	if !exists {
		return fmt.Errorf("node not found")
	}

	var collectionFound *data.CollectionDefinition
	for _, collection := range orchestrator.node.Collections {
		if collection.Name == collectionName {
			collectionFound = &collection
			break
		}
	}

	if collectionFound == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	collectionFound.Fields[field] = fieldType

	//update manifest
	if err := orchestrator.updateNodeManifest(); err != nil {
		return err
	}

	//restart container
	return orchestrator.restartContainer(nodeData.ContainerId)
}

//RemoveField remove a single field from the specified collection.
//Node metadata is also updated.
func (orchestrator *NodeOrchestrator) RemoveField(collectionName string, field string) error {
	nodeData, exists := orchestrator.persistence.GetNode(orchestrator.node.Name)
	if !exists {
		return fmt.Errorf("node not found")
	}

	var collectionFound *data.CollectionDefinition
	for _, collection := range orchestrator.node.Collections {
		if collection.Name == collectionName {
			collectionFound = &collection
			break
		}
	}

	if collectionFound == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	if _, exists := collectionFound.Fields[field]; !exists {
		return fmt.Errorf("field %s not found", field)
	}

	delete(collectionFound.Fields, field)

	//update manifest
	if err := orchestrator.updateNodeManifest(); err != nil {
		return err
	}

	//restart container
	return orchestrator.restartContainer(nodeData.ContainerId)
}

//createContainer creates a new docker container, the container id is returned or error if something went wrong
func (orchestrator *NodeOrchestrator) createContainer() (string, error) {
	ctx := context.Background()

	hostName := fmt.Sprintf("apio-%s.%s", orchestrator.node.Name, orchestrator.nodeDockerConfig.DomainName)

	containerConfig := &container.Config{
		Image: config.ApiDockerImage,
		Env: []string{
			fmt.Sprintf("MONGODB_HOST=%s", orchestrator.nodeDockerConfig.MongoDbHost),
			fmt.Sprintf("MONGODB_NAME=%s", orchestrator.nodeDockerConfig.MongoDbName),
			"DEBUG_MODE=1",
		},
		Labels: proxy.GetTraefikLabels("apio_"+orchestrator.node.Name, hostName, orchestrator.nodeDockerConfig.NetworkName, 80),
	}
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: node.ManifestFile(orchestrator.node.Name),
				Target: "/app/manifest.json",
			},
		},
	}
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			orchestrator.nodeDockerConfig.NetworkName: {},
		},
	}

	result, err := orchestrator.dockerClient.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		networkConfig,
		fmt.Sprintf("%s%s", orchestrator.nodeDockerConfig.ContainerNamePrefix, orchestrator.node.Name),
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

func (orchestrator *NodeOrchestrator) startContainer(containerId string) error {
	return orchestrator.dockerClient.ContainerStart(context.Background(), containerId, types.ContainerStartOptions{})
}

func (orchestrator *NodeOrchestrator) stopContainer(containerId string) error {
	timeout := time.Second * 10
	return orchestrator.dockerClient.ContainerStop(context.Background(), containerId, &timeout)
}

func (orchestrator *NodeOrchestrator) restartContainer(containerId string) error {
	timeout := time.Second * 10
	return orchestrator.dockerClient.ContainerRestart(context.Background(), containerId, &timeout)
}

func (orchestrator *NodeOrchestrator) createNodeManifest() error {
	//using a different method name just in case we need to add more validations depending if it's an update or creation
	return orchestrator.updateNodeManifest()
}

func (orchestrator *NodeOrchestrator) updateNodeManifest() error {
	data, _ := json.MarshalIndent(orchestrator.node.Collections, "", "    ")
	return ioutil.WriteFile(node.ManifestFile(orchestrator.node.Name), data, 0644)
}

func (orchestrator *NodeOrchestrator) ensureNetworkIsCreated() error {
	ctx := context.Background()

	if orchestrator.networkExists() {
		return nil
	}
	_, err := orchestrator.dockerClient.NetworkCreate(ctx, orchestrator.nodeDockerConfig.NetworkName, types.NetworkCreate{})
	if err != nil {
		return err
	}
	return nil
}

func (orchestrator *NodeOrchestrator) networkExists() bool {
	ctx := context.Background()

	args := filters.NewArgs()
	args.Add("name", orchestrator.nodeDockerConfig.NetworkName)
	listResult, err := orchestrator.dockerClient.NetworkList(ctx, types.NetworkListOptions{Filters: args})
	if err != nil {
		return false
	}

	for _, networkResult := range listResult {
		if networkResult.Name == orchestrator.nodeDockerConfig.NetworkName {
			//network found
			return true
		}
	}
	return false
}

func (orchestrator *NodeOrchestrator) fetchImage() error {
	//if _, err := orchestrator.dockerClient.RegistryLogin(ctx, types.AuthConfig{
	//	Username: "...",
	//	Password: "...token",
	//	ServerAddress: "docker.pkg.github.com",
	//}); err != nil {
	//	return err
	//}
	//if _, err := orchestrator.dockerClient.ImagePull(ctx, apiDockerImage, types.ImagePullOptions{
	//	RegistryAuth: "...",
	//}); err != nil {
	//	return err
	//}
	return nil
}

func (orchestrator *NodeOrchestrator) removeContainer(containerId string) error {
	return orchestrator.dockerClient.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{
		Force: true,
	})
}
