package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chrismaher/redsheets/homedir"
	"github.com/chrismaher/redsheets/json"
)

type connection struct {
	Host string
	Port string
	User string
	Name string
}

var (
	cfgFile   string
	dataStore string
	conn      connection
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "redsheets",
	Short: "redsheets provides an interfact for managing Google Sheets -> Redshift mappings",
	Long:  `A CLI for managing mappings between Google Sheets and Redshift tables. Complete documentation is available at https://github.com/chrismaher/redsheets.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.redsheets/config.toml)")
	rootCmd.PersistentFlags().StringVar(&dataStore, "datastore", "", "")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.FullPath(dataStore)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".redsheets" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("config.toml")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&conn); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data = json.Data{Path: "/Users/cmaher/.redsheets/data/redsheets.json"}
	data.Read()
}
