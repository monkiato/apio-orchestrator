package cmd

import (
	"errors"
	"fmt"
	"github.com/monkiato/apio-orchestrator/internal/data"
	"github.com/monkiato/apio-orchestrator/internal/tools"
	"github.com/spf13/cobra"
)

var (
	addCollectionCmd = &cobra.Command{
		Use:   "addCollection [node id] [collection name]",
		Short: "create a new collection definition for the specified node id",
		Long: `AddCollection (api-orchestrator addCollection) will create a new and empty collection definition for the
specified node id, then addField can be used for adding fields to this collection.

Example: apio-orchestrator addCollection my-client-crm authors`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := addCollection(args[0], args[1]); err != nil {
				onError(err)
			}
			fmt.Printf("Collection '%s' added successfully for node id '%s'", args[1], args[0])
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

func addCollection(nodeId string, collectionName string) error {
	nodeOrchestrator := createNodeOrchestrator(nodeId)
	return nodeOrchestrator.AddCollection(data.CollectionDefinition{
		Name: collectionName,
		Fields: map[string]string{},
	})
}
