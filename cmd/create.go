package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create [node id]",
		Short: "create a new Apio node in docker",
		Long: `Create (api-orchestrator create) will create a new Apio node in your local docker
instance. If a manifest path is not provided the node is created with no collections.

Example: apio-orchestrator create my-client-crm`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := createNode(args[0]); err != nil {
				onError(err)
			}
			fmt.Printf("Apio node '%s' created successfully", args[0])
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
	createCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "", "path to manifest file")
}

func createNode(nodeId string) error {
	collections, err := readCollections(manifestPath)
	if err != nil {
		return err
	}

	nodeOrchestrator := createNodeOrchestrator(nodeId)

	_, err = nodeOrchestrator.CreateNode()
	if err != nil {
		return err
	}

	if err := nodeOrchestrator.UpdateCollections(collections); err != nil {
		return err
	}

	return nil
}
