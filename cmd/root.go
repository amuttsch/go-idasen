package cmd

import (
	"fmt"
	"github.com/amuttsch/go-idasen/idasen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "go-idasen",
	Short: "Control your IKEA IDÅSEN desk from command line",
	Long: `A simple and easy to use command line tool to control your height adjustable
	       desk via bluetooth.`,
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
		targetHeight, targetPositionExists := config.Positions[positionName]
		if !targetPositionExists {
			fmt.Printf("Target \"%s\" not found\n", positionName)
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func getConfigFilePath() string {
	path, _ := os.UserConfigDir()

	return path  + "/idasen/idasen.yaml"
}

func initConfig() {
	viper.SetConfigType("yaml")

	viper.SetConfigFile(getConfigFilePath())

	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}
