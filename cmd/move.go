package cmd

import (
	"fmt"
	"github.com/amuttsch/go-idasen/idasen"
	"github.com/spf13/viper"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(moveCmd)
}

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move desk to target height",
	Long:  `Move desk to target height, usage: move 0.7`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := getConfigFilePath()
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Printf("Desk not initialized. Run \"go-idasen init\" first\n")
			return
		}

		if verbose {
			idasen.SetDebug()
		}

		targetHeight, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Printf("Invalid target height \"%s\" given: %s\n", args[0], err)
			return
		}

		config := idasen.Configuration{}
		err = viper.Unmarshal(&config)
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

		err = idasen.MoveToTarget(targetHeight)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
