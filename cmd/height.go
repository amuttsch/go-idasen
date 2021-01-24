package cmd

import (
	"fmt"
	"github.com/amuttsch/go-idasen/idasen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	rootCmd.AddCommand(heightCmd)
}

var heightCmd = &cobra.Command{
	Use:   "height",
	Short: "Show current height",
	Long:  `Connect to the desk and print the current desk height in meters.`,
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

		fmt.Printf("%.4fm\n", h)
	},
}
