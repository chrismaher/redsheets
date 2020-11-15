package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/chrismaher/redsheets/json"
	"github.com/chrismaher/redsheets/redshift"
	"github.com/spf13/cobra"
)

var (
	schema    string
	name      string
	id        string
	sheet     string
	datastore json.DataStore
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a GoogleSheets-to-Redshift mapping",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		table := json.NewMap(id, sheet, schema, name)
		ln, err := datastore.Add(table)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added map %d\n", ln)
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
				fmt.Println(err)
				return
			}

			err = datastore.Delete(id)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Deleted map %d\n", id)
		}
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mappings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		datastore.Print()
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
		for _, v := range args {
			key, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			table, err := datastore.Get(key)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = service.Authorize()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			sheet_contents, err := service.GetRange(table.SheetID, table.SheetName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			db := redshift.Client{Connection: &connect}
			db.Connect()
			defer db.DB.Close()

			if err = db.Replace(table.Schema, table.Name, sheet_contents[1:]); err != nil {
				fmt.Println(err)
				os.Exit(1)
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
