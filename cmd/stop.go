package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	stopCmd = &cobra.Command{
		Use:   "stop [node id]",
		Short: "stop an existing docker container for the specified Apio node",
		Long: `Stop (api-orchestrator stop) will stop an existing Apio node in your local docker
instance.

Example: apio-orchestrator stop my-client-crm`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := stopNode(args[0]); err != nil {
				onError(err)
			}
			fmt.Printf("Apio node '%s' stoped successfully\n", args[0])
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

func stopNode(nodeId string) error {
	nodeOrchestrator := createNodeOrchestrator(nodeId)
	return nodeOrchestrator.StopNode()
}
