package json

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func TabPrinter(inputs []Table) {
	str := "%v\t%v\t%v\t%v\t\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, str, "Sheet ID", "Sheet Name", "Schema", "Name")
	fmt.Fprintf(tw, str, "--------", "----------", "------", "----")
	for _, val := range inputs {
		fmt.Fprintf(tw, str, val.SheetID, val.SheetName, val.Schema, val.Name)
	}
	tw.Flush()
}
