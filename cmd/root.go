package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chrismaher/redsheets/google"
	"github.com/chrismaher/redsheets/homedir"
	"github.com/chrismaher/redsheets/json"
	"github.com/chrismaher/redsheets/redshift"
)

var (
	configFile string
	dataFile   string
	secretFile string
	tokenFile  string
	service    google.Client
	connect    redshift.Connection
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

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.redsheets/config.toml)")
	rootCmd.PersistentFlags().StringVar(&secretFile, "secret", "", "")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token", "", "")
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	// Find application directory; create if it doesn't exist?
	projectDir, err := homedir.AbsPath(".redsheets")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if configFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(configFile)
	} else {
		// Search for config in ~/.redsheets directory
		viper.AddConfigPath(projectDir)
		viper.SetConfigName("config.toml")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv() // read environment variables that match

	// If a config file is found, read it
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.UnmarshalKey("database", &connect); err != nil {
		// need to test this
		fmt.Println(err)
		os.Exit(1)
	}

	if dataFile == "" {
		dataFile, err = homedir.AbsPath(projectDir, "data.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	datastore = json.DataStore{Path: dataFile}
	datastore.Read()

	if secretFile == "" {
		secretFile, err = homedir.AbsPath(projectDir, viper.GetString("client.secret"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if tokenFile == "" {
		tokenFile, err = homedir.AbsPath(projectDir, viper.GetString("client.token"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	service, err = google.New(secretFile, tokenFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
