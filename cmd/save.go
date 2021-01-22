package cmd

import (
	"fmt"
	"github.com/amuttsch/go-idasen/idasen"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(saveCmd)
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save the current position",
	Long:  `Save the current position under the given name. It can be later used to move to this position using go-idasen [alias]`,
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

		idasen, err := idasen.New(config)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer idasen.Close()

		h, err := idasen.HeightInMeters()
		if err != nil {
			fmt.Println(err)
			return
		}

		positionName := args[0]
		config.Positions[positionName] = h

		viper.Set("mac_address", config.MacAddress)
		viper.Set("positions", config.Positions)

		fmt.Printf("Saved position %.4fm as %s\n", h, positionName)

		_ = viper.WriteConfigAs(configFile)
	},
}
