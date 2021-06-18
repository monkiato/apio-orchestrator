package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update [node id]",
		Short: "update an existing docker container for the specified Apio node",
		Long: `Update (api-orchestrator update) will update an existing Apio node in your local docker
instance.

Example: apio-orchestrator update my-client-crm`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := updateNode(args[0]); err != nil {
				onError(err)
			}
			fmt.Printf("Apio node '%s' updateed successfully\n", args[0])
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
	updateCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "", "path to manifest file")
	updateCmd.MarkFlagRequired("manifest")
}

func updateNode(nodeId string) error {
	collections, err := readCollections(manifestPath)
	if err != nil {
		return err
	}
	nodeOrchestrator := createNodeOrchestrator(nodeId)
	return nodeOrchestrator.UpdateCollections(collections)
}
