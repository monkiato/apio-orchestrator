package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove [node id]",
		Short: "remove an existing docker container for the specified Apio node",
		Long: `Remove (api-orchestrator remove) will remove an existing Apio node in your local docker
instance.

Example: apio-orchestrator remove my-client-crm`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := removeNode(args[0]); err != nil {
				onError(err)
			}
			fmt.Printf("Apio node '%s' removed successfully", args[0])
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
		Aliases: []string{"rm"},
	}
)

func init() {
}

func removeNode(nodeId string) error {
	nodeOrchestrator := createNodeOrchestrator(nodeId)
	return nodeOrchestrator.RemoveNode()
}
