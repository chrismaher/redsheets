package cmd

import (
	"fmt"

	"github.com/chrismaher/redsheets/google"
	"github.com/chrismaher/redsheets/json"
	"github.com/chrismaher/redsheets/redshift"
	"github.com/spf13/cobra"
	"log"
)

var (
	schema    string
	name      string
	sheetID   string
	sheetName string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a JSON datastore of GoogleSheets-to-Redshift mappings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := json.Init(); err != nil {
			fmt.Println(err)
		}
	},
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a GoogleSheets-to-Redshift mapping",
	Long:  ``,
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
	Short: "Delete a GoogleSheets-to-Redshift mapping",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mappings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := json.List(); err != nil {
			fmt.Println(err)
		}
	},
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a GoogleSheets-to-Redshift mapping",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
	},
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a GoogleSheets-to-Redshift refresh for a given sheet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		js, err := json.Read()
		if err != nil {
			log.Panic(err)
		}
		var table json.Table

		for _, tab := range js {
			if tab.Schema == schema && tab.Name == name {
				table = tab
			}
		}

		service := google.Service{}

		err = service.Authorize()
		if err != nil {
			log.Panic(err)
		}

		sheet_contents, err := service.GetRange(table.SheetID, table.SheetName)
		if err != nil {
			log.Panic(err)
		}

		db := redshift.Client{}
		db.Connect()
		defer db.DB.Close()

		if err = db.Replace(table.Schema, table.Name, sheet_contents[1:]); err != nil {
			log.Panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&schema, "schema", "", "The target tables's schema (required)")
	runCmd.Flags().StringVar(&name, "table", "", "The target tables's name (required)")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)

	addCmd.Flags().StringVar(&sheetID, "sheet_id", "", "AWS region (required)")
	addCmd.Flags().StringVar(&sheetName, "sheet_name", "", "AWS region (required)")
	addCmd.Flags().StringVar(&schema, "schema", "", "AWS region (required)")
	addCmd.Flags().StringVar(&name, "table", "", "AWS region (required)")
}
