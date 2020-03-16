package cmd

import (
	"github.com/chrismaher/redsheets/google"
	"github.com/chrismaher/redsheets/json"
	"github.com/chrismaher/redsheets/redshift"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a GoogleSheets-to-Redshift refresh for a given sheet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		js := json.Read()
		var table json.Table

		for _, tab := range js {
			if tab.Schema == schema && tab.Name == name {
				table = tab
			}
		}

		service := google.Service{}
		service.Authorize()

		spreadsheetId := table.SheetID
		readRange := table.Name

		sheet_contents := service.GetRange(spreadsheetId, readRange)

		db := redshift.Client{}
		db.Connect()
		defer db.DB.Close()

		db.Replace(table.Schema, table.Name, sheet_contents[1:])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&schema, "schema", "", "The target tables's schema (required)")
	runCmd.Flags().StringVar(&name, "table", "", "The target tables's name (required)")
}
