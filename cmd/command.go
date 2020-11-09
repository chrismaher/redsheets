package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chrismaher/redsheets/google"
	"github.com/chrismaher/redsheets/json"
	"github.com/chrismaher/redsheets/redshift"
	"github.com/spf13/cobra"
)

var (
	schema string
	name   string
	id     string
	sheet  string
	data   json.Data
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a GoogleSheets-to-Redshift mapping",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		table := json.Table{id, sheet, schema, name}
		key, err := data.Add(table)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Created table %d\n", key)
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a GoogleSheets-to-Redshift mapping",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			id, err := strconv.Atoi(v)
			if err != nil {
				log.Println(err)
				return
			}
			data.Delete(id)
			fmt.Printf("Deleted table %d\n", id)
		}
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mappings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		data.List()
	},
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a GoogleSheets-to-Redshift mapping",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("update called")
	},
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a GoogleSheets-to-Redshift refresh for a given sheet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			key, err := strconv.Atoi(v)
			if err != nil {
				log.Panic(err)
			}

			table, err := data.Get(key)
			if err != nil {
				log.Panic(err)
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

			db := redshift.Client{Connection: &conn}
			db.Connect()
			defer db.DB.Close()

			if err = db.Replace(table.Schema, table.Name, sheet_contents[1:]); err != nil {
				log.Panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(updateCmd)

	addCmd.Flags().StringVar(&id, "sheet_id", "", "AWS region (required)")
	addCmd.Flags().StringVar(&sheet, "sheet_name", "", "AWS region (required)")
	addCmd.Flags().StringVar(&schema, "schema", "", "AWS region (required)")
	addCmd.Flags().StringVar(&name, "table", "", "AWS region (required)")
}
