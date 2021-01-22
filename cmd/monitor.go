package cmd

import (
	"fmt"
	"github.com/amuttsch/go-idasen/idasen"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(monitorCmd)
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor height changes",
	Long:  `Continuously monitor the height changes of the desk`,
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

		var wg sync.WaitGroup
		wg.Add(1)

		ch := make(chan os.Signal,1 )
		signal.Notify(ch, os.Interrupt, os.Kill)

		go func() {
			var previousHeight = 0.0

			for true {
				h, err := idasen.HeightInMeters()
				if err != nil {
					fmt.Println(err)
					return
				}

				if h != previousHeight {
					fmt.Printf("%.4fm\n", h)
					previousHeight = h
				}

				select {
				case <-ch:
					wg.Done()
					return
				default:

				}
			}
		}()

		wg.Wait()
	},
}
