package cmd

import (
	"fmt"

	"github.com/chrismaher/redsheets/json"
	"github.com/spf13/cobra"
)

var (
	schema    string
	name      string
	sheetID   string
	sheetName string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := json.Init(); err != nil {
			fmt.Println(err)
		}
	},
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		table := json.Table{sheetID, sheetName, schema, name}
		if err := json.Add(table); err != nil {
			fmt.Println(err)
		}
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := json.List(); err != nil {
			fmt.Println(err)
		}
	},
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	addCmd.Flags().StringVar(&sheetID, "sheet_id", "", "AWS region (required)")
	addCmd.Flags().StringVar(&sheetName, "sheet_name", "", "AWS region (required)")
	addCmd.Flags().StringVar(&schema, "schema", "", "AWS region (required)")
	addCmd.Flags().StringVar(&name, "table", "", "AWS region (required)")

	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(addCmd)
	configCmd.AddCommand(deleteCmd)
	configCmd.AddCommand(listCmd)
	configCmd.AddCommand(updateCmd)
}