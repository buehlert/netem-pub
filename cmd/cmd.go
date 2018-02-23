package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thomas-fossati/netem-pub/cmd/config"
)

var Verbose bool
var cfg *config.Config
var cfgFile string

var NetemPubCmd = &cobra.Command{
	Use:   "TODO",
	Short: "TODO high-level...",
	Long: `TODO details...
			Complete documentation is available at http://...`,
	Run: func(cmd *cobra.Command, args []string) {
		NetemPubSvc(cfg)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	NetemPubCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is TODO)")
	NetemPubCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func initConfig() {
	cfg = config.NewConfig()
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME/.netem-pub")
		viper.AddConfigPath("/etc/netem-pub")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Println("Can't unmarshal config:", err)
		os.Exit(1)
	}
}