package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start [node id]",
		Short: "start an existing docker container for the specified Apio node",
		Long: `Start (api-orchestrator start) will start an existing Apio node in your local docker
instance.

Example: apio-orchestrator start my-client-crm`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := startNode(args[0]); err != nil {
				onError(err)
			}
			fmt.Printf("Apio node '%s' started successfully", args[0])
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

func startNode(nodeId string) error {
	nodeOrchestrator := createNodeOrchestrator(nodeId)
	return nodeOrchestrator.StartNode()
}
