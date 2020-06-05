package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/monkiato/apio-orchestrator/internal/data"
	"github.com/monkiato/apio-orchestrator/internal/tools"
	"github.com/monkiato/apio-orchestrator/pkg/persistence"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

const (
	validNodeIdChars = "abcdefghijklmnopqrstuvwxyz1234567890-_"
)

var (
	persistenceConnection persistence.Connection
	manifestPath          string

	rootCmd = &cobra.Command{
		Use:   "apio-orchestrator [OPTIONS]",
		Short: "Apio orchestrator CLI for Apio node management.",
		Long: `Apio orchestrator CLI is a management tool for Apio nodes running in docker with the ability
to handle new or existing Apio containers and modify their collections`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	persistenceConnection = persistence.NewFileConnection()

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(addCollectionCmd)
	rootCmd.AddCommand(removeCollectionCmd)
	rootCmd.AddCommand(addFieldCmd)
	rootCmd.AddCommand(removeFieldCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(inspectCmd)
}

func onError(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {

}

func isValidNodeId(nodeId string) bool {
	return tools.IsValidFormat(nodeId, validNodeIdChars)
}

func readCollections(manifestPath string) ([]data.CollectionDefinition, error) {
	collections := []data.CollectionDefinition{}

	if manifestPath != "" {
		data, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &collections); err != nil {
			return nil, err
		}
	}
	return collections, nil
}
