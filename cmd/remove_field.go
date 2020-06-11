package cmd

import (
	"errors"
	"fmt"
	"github.com/monkiato/apio-orchestrator/internal/tools"
	"github.com/spf13/cobra"
)

var (
	removeFieldCmd = &cobra.Command{
		Use:   "removeField [node id] [collection name]",
		Short: "remove an existing field definition for the specified node id and collection",
		Long: `removeField (api-orchestrator removeField) will delete an existing field definition for the
specified node id and collection.

Example: apio-orchestrator removeField my-client-crm authors name`,
		Run: func(cmd *cobra.Command, args []string) {
			nodeId := args[0]
			collection := args[1]
			field := args[2]
			if err := removeField(nodeId, collection, field); err != nil {
				onError(err)
			}
			fmt.Printf("Field '%s' removed successfully for node id '%s' and collection '%s", field, nodeId, collection)
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
			if !tools.IsValidFormat(args[2], validNodeIdChars) {
				return errors.New("field name format is invalid. Expected alphanumeric value and '-' or '_' as word separator")
			}
			return nil
		},
	}
)

func init() {
}

func removeField(nodeId string, collectionName string, field string) error {
	nodeOrchestrator := createNodeOrchestrator(nodeId)
	return nodeOrchestrator.RemoveField(collectionName, field)
}
