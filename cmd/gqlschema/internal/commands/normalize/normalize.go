package normalize

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
		Use:   "normalize",
		Short: "Normalize graphql schema file",
		Long: "Normalize graphql schema: standardize structure, sort fields alphabetically\n" +
			"If the file is not provided, it will be read from stdin",
		Run: normalize,
	}
	c.Flags().StringVarP(&f.schemaFile, "file", "f", "", "Schema file to normalize")
	c.Flags().StringVarP(&f.outputFile, "output", "o", "", "Target file to output to")
	return &c
}

func normalize(cmd *cobra.Command, _ []string) {
	if (f.schemaFile == "") && (term.IsTerminal(os.Stdin.Fd())) {
		fmt.Fprintf(os.Stderr, "No input file provided (use -f to specify a file)\n\n")
		if err := cmd.Usage(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
	schema, err := gqls.LoadSDL(f.schemaFile)
	cobra.CheckErr(err)

	gqls.Normalize(schema)

	err = gqls.OutputFile(schema, f.outputFile)
	cobra.CheckErr(err)
}
