package convert

import (
	"fmt"
	"os"

	"github.com/catonetworks/cato-go-sdk/gqls"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

type cmdFlags struct {
	schemaFile string
	outputFile string
}

var f cmdFlags

func Cmd() *cobra.Command {
	var c = cobra.Command{
		Use:   "convert",
		Short: "Convert graphql schema introspection JSON into and SDL format",
		Run:   convert,
	}
	c.Flags().StringVarP(&f.schemaFile, "file", "f", "", "JSON introspection file to convert")
	c.Flags().StringVarP(&f.outputFile, "output", "o", "", "Target SDL file to output to")
	return &c
}

func convert(cmd *cobra.Command, _ []string) {
	if (f.schemaFile == "") && (term.IsTerminal(os.Stdin.Fd())) {
		fmt.Fprintf(os.Stderr, "No input file provided (use -f to specify a file)\n\n")
		if err := cmd.Usage(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
	schema, err := gqls.LoadJSON(f.schemaFile)
	cobra.CheckErr(err)

	gqls.Normalize(schema)

	err = gqls.OutputFile(schema, f.outputFile)
	cobra.CheckErr(err)
}
