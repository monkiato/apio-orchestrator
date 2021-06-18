package cmd

import (
	"errors"
	"fmt"
	"github.com/monkiato/apio-orchestrator/internal/tools"
	"github.com/spf13/cobra"
)

var (
	addFieldCmd = &cobra.Command{
		Use:   "addField [node id] [collection name]",
		Short: "create a new field definition for the specified node id and collection",
		Long: `addField (api-orchestrator addField) will create a new field definition for the
specified node id and collection.

Available type of fields: 'string', 'bool', 'float'

Example: apio-orchestrator addField my-client-crm authors name string`,
		Run: func(cmd *cobra.Command, args []string) {
			nodeId := args[0]
			collection := args[1]
			field := args[2]
			fieldType := args[3]
			if err := addField(nodeId, collection, field, fieldType); err != nil {
				onError(err)
			}
			fmt.Printf("New field '%s' (%s) added successfully for node id '%s' and collection '%s\n", field, fieldType, nodeId, collection)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("node id is required")
			}
			if len(args) < 2 {
				return errors.New("collection name is required")
			}
			if len(args) < 3 {
				return errors.New("field name is required")
			}
			if len(args) < 4 {
				return errors.New("field type is required, expected 'string', 'bool' or 'float'")
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
			if _, found := tools.Find([]string{"string", "bool", "float"}, args[3]); !found {
				return errors.New("unexpected field type. Valid formats are: \"string\", \"bool\", \"float\"")
			}
			return nil
		},
	}
)

func init() {
}

func addField(nodeId string, collectionName string, field string, fieldType string) error {
	nodeOrchestrator := createNodeOrchestrator(nodeId)
	return nodeOrchestrator.AddField(collectionName, field, fieldType)
}
