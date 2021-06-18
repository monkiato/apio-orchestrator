package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/monkiato/apio-orchestrator/internal/data"
	"github.com/monkiato/apio-orchestrator/internal/tools"
	"github.com/monkiato/apio-orchestrator/pkg/config"
	"github.com/monkiato/apio-orchestrator/pkg/node"
	"github.com/monkiato/apio-orchestrator/pkg/orchestrator"
	"github.com/monkiato/apio-orchestrator/pkg/persistence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

const (
	validNodeIdChars = "abcdefghijklmnopqrstuvwxyz1234567890-_"
)

var (
	persistenceConnection persistence.Connection
	nodeDockerConfig      *node.DockerConfig
	cfgFile               string
	manifestPath          string

	rootCmd = &cobra.Command{
		Use:   "apio-orchestrator [OPTIONS]",
		Short: "Apio orchestrator CLI for Apio node management.",
		Long: `Apio orchestrator CLI is a management tool for Apio nodes running in docker with the ability
to handle new or existing Apio containers and modify their collections`,
	}
)

//Execute run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault("domain_name", config.DefaultDomainName)
	viper.SetDefault("network_name", config.DefaultNetworkName)
	viper.SetDefault("apio_node_prefix", config.DefaultNodePrefix)
	viper.SetDefault("config_path", config.DefaultConfigPath)
	viper.SetDefault("mongodb.host", config.DefaultMongoDbHost)
	viper.SetDefault("mongodb.name", config.DefaultMongoDbName)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/apio-orchestrator/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.apio-orchestrator") // call multiple times to add many search paths
	viper.AddConfigPath(".")                        // optionally look for config in the working directory

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file absolute path")

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

func isDir(path string) (bool, error) {
	fi, err := os.Stat(path)
    if err != nil {
        return false, fmt.Errorf("Unable to read file stats")
    }
    mode := fi.Mode()
	return mode.IsDir(), nil
}

func initConfig() {
	if cfgFile != "" {
		if r, err := isDir(cfgFile); err != nil || r {
			fmt.Printf("error: invalid config file '%s'\n", cfgFile)
			os.Exit(1)
		}
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}

	configPath := viper.GetString("config_path")
	persistenceConnection = persistence.NewFileConnection(configPath)
	node.SetRootConfigPath(configPath)

	nodeDockerConfig = &node.DockerConfig{
		NetworkName:         viper.GetString("network_name"),
		DomainName:          viper.GetString("domain_name"),
		ContainerNamePrefix: viper.GetString("apio_node_prefix"),
		MongoDbHost:         viper.GetString("mongodb.host"),
		MongoDbName:         viper.GetString("mongodb.name"),
	}
}

func isValidNodeId(nodeId string) bool {
	return tools.IsValidFormat(nodeId, validNodeIdChars)
}

func createNodeOrchestrator(nodeId string) *orchestrator.NodeOrchestrator {
	nodeOrchestrator, err := orchestrator.NewNodeOrchestrator(nodeId, nodeDockerConfig, persistenceConnection)
	if err != nil {
		onError(err)
	}
	return nodeOrchestrator
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
