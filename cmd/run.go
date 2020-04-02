package cmd

import (
	"github.com/chrismaher/redsheets/google"
	"github.com/chrismaher/redsheets/json"
	"github.com/chrismaher/redsheets/redshift"
	"github.com/spf13/cobra"
	"log"
)

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

		db.Replace(table.Schema, table.Name, sheet_contents[1:])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&schema, "schema", "", "The target tables's schema (required)")
	runCmd.Flags().StringVar(&name, "table", "", "The target tables's name (required)")
}
