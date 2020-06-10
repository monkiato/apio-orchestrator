package cmd

import (
	"fmt"
	"github.com/monkiato/apio-orchestrator/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

var (
	lsCmd = &cobra.Command{
		Use:   "ls [node id]",
		Short: "list all nodes",
		Long: `ls (api-orchestrator ls) will list all existing Apio node in your local docker
instance.

Example: apio-orchestrator ls`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := listNodes(); err != nil {
				onError(err)
			}
		},
	}
)

func init() {
}

func listNodes() error {
	fmt.Println("NODES")
	files, err := ioutil.ReadDir(viper.GetString("config_path") + config.NodeFolder)
	if err != nil {
		return err
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			fmt.Println(fileInfo.Name())
		}
	}
	return nil
}
