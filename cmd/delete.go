package cmd

import (
	"fmt"
	"github.com/amuttsch/go-idasen/idasen"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a saved position",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := getConfigFilePath()
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Printf("Desk not initialized. Run \"go-idasen init\" first\n")
			return
		}

		if verbose {
			idasen.SetDebug()
		}

		config := idasen.Configuration{}
		err := viper.Unmarshal(&config)
		if err != nil {
			fmt.Printf("Could not read configuration file \"%s\": %s\n", configFile, err)
			return
		}

		positionName := args[0]
		delete(config.Positions, positionName)

		viper.Set("mac_address", config.MacAddress)
		viper.Set("positions", config.Positions)

		fmt.Printf("Removed position \"%s\"\n", positionName)

		_ = viper.WriteConfigAs(configFile)
	},
}
