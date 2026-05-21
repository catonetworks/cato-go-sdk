package fetch

import (
	"bytes"
	"os"
	"time"

	"github.com/catonetworks/cato-go-sdk/gqls"
	"github.com/spf13/cobra"
)

type cmdFlags struct {
	introspectURL string
	outputFile    string
	convert       bool
	all           bool
	token         string
}

var f cmdFlags

const timeout = time.Second * 10

func Cmd() *cobra.Command {
	var c = cobra.Command{
		Use:   "fetch",
		Short: "Fetch schema introspection from a given grqaphql server",
		Run:   fetch,
	}
	c.Flags().StringVarP(&f.introspectURL, "url", "u", "", "URL of the graphql server (var: $CATO_ENDPOINT)")
	c.Flags().StringVarP(&f.outputFile, "output", "o", "", "Target file to output to")
	c.Flags().StringVarP(&f.token, "token", "t", "", "Authentication token (var: $CATO_TOKEN)")
	c.Flags().BoolVarP(&f.convert, "convert", "c", false, "Convert to an SDL format (default is JSON)")
	c.Flags().BoolVarP(&f.all, "all", "a", false, "Fetch all types, including the undocumented ones")
	return &c
}

func fetch(_ *cobra.Command, _ []string) {
	// URL
	if f.introspectURL == "" {
		f.introspectURL = os.Getenv("CATO_ENDPOINT")
	}
	if f.introspectURL == "" {
		cobra.CheckErr("Introspection URL not provided. Please provide it with the --url flag or set the CATO_ENDPOINT environment variable.")
	}
	// Token
	if f.token == "" {
		f.token = os.Getenv("CATO_TOKEN")
	}
	if f.token == "" {
		cobra.CheckErr("Introspection URL not provided. Please provide it with the --url flag or set the CATO_TOKEN environment variable.")
	}

	// fetch data
	client := gqls.NewHTTPClient(timeout)
	introspectData, err := client.GetIntrospection(f.introspectURL, f.token, f.all)
	cobra.CheckErr(err)

	// print JSON data to stdout
	if !f.convert {
		if f.outputFile != "" {
			cobra.CheckErr(os.WriteFile(f.outputFile, introspectData, 0o600))
			return
		}
		_, err := os.Stdout.Write(introspectData)
		cobra.CheckErr(err)
		return
	}

	// parse introspection, converto to SDL
	introReader := bytes.NewBuffer(introspectData)
	schema, err := gqls.ParseJSONSchema(introReader, f.introspectURL)
	cobra.CheckErr(err)
	gqls.Normalize(schema)
	// print schema
	err = gqls.OutputFile(schema, f.outputFile)
	cobra.CheckErr(err)
}
