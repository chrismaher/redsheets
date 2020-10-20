package json

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func TabPrinter(inputs Tables) {
	str := "%v\t%v\t%v\t%v\t%v\t\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, str, "ID", "Sheet ID", "Sheet Name", "Schema", "Name")
	fmt.Fprintf(tw, str, "--", "--------", "----------", "------", "----")

	var v Table
	for i := 0; i < len(inputs); i++ {
		v = inputs[i]
		fmt.Fprintf(tw, str, i, v.SheetID, v.SheetName, v.Schema, v.Name)
	}

	tw.Flush()
}
