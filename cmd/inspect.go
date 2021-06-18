package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/monkiato/apio-orchestrator/pkg/node"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var (
	inspectCmd = &cobra.Command{
		Use:   "inspect [node id]",
		Short: "inspect node metadata",
		Long: `inspect (api-orchestrator inspect) will provide metadata info for an existing Apio node in your local docker
instance.

Example: apio-orchestrator inspect my-client-crm`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := inspectNode(args[0]); err != nil {
				onError(err)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("node id is required")
			}
			if !isValidNodeId(args[0]) {
				return errors.New("node id format is invalid. Expected alphanumeric value and '-' or '_' as word separator")
			}
			return nil
		},
	}
)

func init() {
}

func inspectNode(nodeId string) error {
	metadata, err := ioutil.ReadFile(node.MetadataFile(nodeId))
	if err != nil {
		return fmt.Errorf("metadata not found for node ID '%s'", nodeId)
	}
	manifestData, err := ioutil.ReadFile(node.ManifestFile(nodeId))
	if err != nil {
		return fmt.Errorf("manifest not found for node ID '%s'", nodeId)
	}

	//marshall and unmarshall again for better indentation
	var data interface{}
	err = json.Unmarshal([]byte(fmt.Sprintf("{\"metadata\":%s,\"manifest\":%s}", metadata, manifestData)), &data)
	if err != nil {
		return err
	}
	jsonWithIndentation, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to parse JSON manifest and metadata for node ID '%s'", nodeId)
	}

	fmt.Println(string(jsonWithIndentation))
	return nil
}
