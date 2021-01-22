package cmd

import (
	"fmt"
	"github.com/amuttsch/go-idasen/idasen"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var forceInit bool

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().BoolVarP(&forceInit, "force", "f", false, "force init")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize desk",
	Long:  `Search for the desk and save the mac address to the configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := getConfigFilePath()
		if _, err := os.Stat(configFile); !os.IsNotExist(err) && !forceInit {
			fmt.Printf("Desk already initialized. Using config file \"%s\"\n", viper.ConfigFileUsed())
			return
		}

		if verbose {
			idasen.SetDebug()
		}

		desk, err := idasen.DiscoverDesk()
		if err != nil {
			fmt.Println(err)
			return
		}
		
		config := &idasen.Configuration{
			MacAddress: desk.Address,
			Positions: map[string]float64{
				"sit": 0.75,
				"stand": 1.11,
			},
		}
		viper.Set("mac_address", config.MacAddress)
		viper.Set("positions", config.Positions)

		_ = viper.WriteConfigAs(configFile)
		fmt.Printf("Written configuration to %s\n", configFile)
	},
}
