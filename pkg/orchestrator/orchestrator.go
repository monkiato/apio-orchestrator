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
	"monkiato/apio-orchestrator/internal/tools/proxy"
	"monkiato/apio-orchestrator/pkg/node"
	"os"
)

const (
	rootNodeFolder = "/Users/rodrigo/apio/nodes/"
	apiDockerImage = "docker.pkg.github.com/monkiato/apio/apio:0.5-alpha"
	domainName     = "monkiato.com"
	networkName    = "apio-network"
	nodePrefix     = "apio-"
)

type NodeOrchestrator struct {
	node         *node.Metadata
	dockerClient *client.Client
}

func NewNodeOrchestrator(node *node.Metadata) (*NodeOrchestrator, error) {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &NodeOrchestrator{
		node:         node,
		dockerClient: dockerClient,
	}, nil
}

func (orchestrator *NodeOrchestrator) CreateNode() error {
	if err := orchestrator.createNodeFolder(); err != nil {
		return err
	}
	if err := orchestrator.createNodeManifest(); err != nil {
		return err
	}
	if err := orchestrator.ensureNetworkIsCreated(); err != nil {
		return err
	}
	if err := orchestrator.deployNode(); err != nil {
		return err
	}
	return nil
}

func (orchestrator *NodeOrchestrator) manifestFile() string {
	return fmt.Sprintf("%s/manifest.json", orchestrator.nodeFolder())
}

func (orchestrator *NodeOrchestrator) nodeFolder() string {
	return fmt.Sprintf("%s%s", rootNodeFolder, orchestrator.node.Name)
}

func (orchestrator *NodeOrchestrator) deployNode() error {
	ctx := context.Background()

	hostName := fmt.Sprintf("apio-%s.%s", orchestrator.node.Name, domainName)

	config := &container.Config{
		Image: apiDockerImage,
		Env: []string{
			"MONGODB_HOST=mongodb:27017",
			"MONGODB_NAME=apio",
			"DEBUG_MODE=1",
		},
		Labels: proxy.GetTraefikLabels("apio_" + orchestrator.node.Name, hostName, networkName, 80),
	}
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: orchestrator.manifestFile(),
				Target: "/app/manifest.json",
			},
		},
	}
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {},
		},
	}

	orchestrator.fetchImage()

	result, err := orchestrator.dockerClient.ContainerCreate(
		ctx,
		config,
		hostConfig,
		networkConfig,
		fmt.Sprintf("%s%s", nodePrefix, orchestrator.node.Name),
	)
	if err != nil {
		return err
	}
	return orchestrator.dockerClient.ContainerStart(ctx, result.ID, types.ContainerStartOptions{})
}

func (orchestrator *NodeOrchestrator) createNodeManifest() error {
	data, err := json.Marshal(orchestrator.node.Collections)
	if err != nil {
		return err
	}

	file, err := os.Create(orchestrator.manifestFile())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (orchestrator *NodeOrchestrator) createNodeFolder() error {
	return os.MkdirAll(orchestrator.nodeFolder(), 0777)
}

func (orchestrator *NodeOrchestrator) ensureNetworkIsCreated() error {
	ctx := context.Background()

	if orchestrator.networkExists() {
		return nil
	}
	_, err := orchestrator.dockerClient.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return err
	}
	return nil
}

func (orchestrator *NodeOrchestrator) networkExists() bool {
	ctx := context.Background()

	args := filters.NewArgs()
	args.Add("name", networkName)
	listResult, err := orchestrator.dockerClient.NetworkList(ctx, types.NetworkListOptions{Filters: args})
	if err != nil {
		return false
	}

	for _, networkResult := range listResult {
		if networkResult.Name == networkName {
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
