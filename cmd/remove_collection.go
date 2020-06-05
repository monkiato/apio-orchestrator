package cmd

import (
	"errors"
	"fmt"
	"github.com/monkiato/apio-orchestrator/internal/tools"
	"github.com/monkiato/apio-orchestrator/pkg/orchestrator"
	"github.com/spf13/cobra"
)

var (
	removeCollectionCmd = &cobra.Command{
		Use:   "removeCollection [node id] [collection name]",
		Short: "remove an existing collection definition for the specified node id",
		Long: `removeCollection (api-orchestrator removeCollection) will remove an existing collection definition for the
specified node id, all field definitions will be removed too.

Example: apio-orchestrator removeCollection my-client-crm authors`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := removeCollection(args[0], args[1]); err != nil {
				onError(err)
			}
			fmt.Printf("Collection '%s' removed successfully for node id '%s'", args[1], args[0])
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("node id is required")
			}
			if len(args) < 2 {
				return errors.New("collection name is required")
			}
			if !isValidNodeId(args[0]) {
				return errors.New("node id format is invalid. Expected alphanumeric value and '-' or '_' as word separator")
			}
			if !tools.IsValidFormat(args[1], validNodeIdChars) {
				return errors.New("collection name format is invalid. Expected alphanumeric value and '-' or '_' as word separator")
			}
			return nil
		},
	}
)

func init() {
}

func removeCollection(nodeId string, collectionName string) error {
	nodeOrchestrator, _ := orchestrator.NewNodeOrchestrator(nodeId, persistenceConnection)
	return nodeOrchestrator.RemoveCollection(collectionName)
}
