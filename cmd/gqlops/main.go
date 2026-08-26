package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/catonetworks/cato-go-sdk/internal/gqlops"
)

const expectedOperations = 544

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(arguments []string, stdout, stderr io.Writer) error {
	if len(arguments) == 0 {
		printUsage(stderr)
		return flag.ErrHelp
	}

	switch arguments[0] {
	case "import":
		return runImport(arguments[1:], stdout, stderr)
	case "check":
		return runCheck(arguments[1:], stdout, stderr)
	case "help", "-h", "--help":
		printUsage(stdout)
		return nil
	default:
		printUsage(stderr)
		return fmt.Errorf("unknown command %q", arguments[0])
	}
}

func runImport(arguments []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("gqlops import", flag.ContinueOnError)
	flags.SetOutput(stderr)

	cliRoot := flags.String("cli-root", "../cato-cli", "path to cato-cli checkout")
	sdkRoot := flags.String("sdk-root", ".", "path to cato-go-sdk checkout")
	expected := flags.Int("expected", expectedOperations, "expected canonical operation count")
	dryRun := flags.Bool("dry-run", false, "validate and report without writing")
	if err := flags.Parse(arguments); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("gqlops import accepts flags only")
	}

	result, err := gqlops.Import(gqlops.ImportConfig{
		CLIRoot:  *cliRoot,
		SDKRoot:  *sdkRoot,
		Expected: *expected,
		DryRun:   *dryRun,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(
		stdout,
		"canonical=%d imported=%d mapped=%d sdk_only=%d\n",
		result.Canonical,
		result.Imported,
		result.Mapped,
		result.SDKOnly,
	)
	return nil
}

func runCheck(arguments []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("gqlops check", flag.ContinueOnError)
	flags.SetOutput(stderr)

	sdkRoot := flags.String("sdk-root", ".", "path to cato-go-sdk checkout")
	expected := flags.Int("expected", expectedOperations, "expected canonical operation count")
	if err := flags.Parse(arguments); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("gqlops check accepts flags only")
	}

	if err := gqlops.Check(gqlops.CheckConfig{SDKRoot: *sdkRoot, Expected: *expected}); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "validated %d canonical operations\n", *expected)
	return nil
}

func printUsage(output io.Writer) {
	fmt.Fprintln(output, "Usage: gqlops <import|check> [flags]")
}
