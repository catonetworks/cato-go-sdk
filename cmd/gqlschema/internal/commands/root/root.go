package root

import (
	"os"

	"github.com/catonetworks/cato-go-sdk/cmd/gqlschema/internal/commands/convert"
	"github.com/catonetworks/cato-go-sdk/cmd/gqlschema/internal/commands/fetch"
	"github.com/catonetworks/cato-go-sdk/cmd/gqlschema/internal/commands/normalize"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// EnvVarPrefix prefix of the environment variables
	EnvVarPrefix = "GQLSCHEMA"
)

var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gqlschema",
	Short: "gqlschema manages your GraphQL schema",
	Long: `gqlschema is a command line tool for managing your GraphQL schema.
It fetches, validates an normalizes it.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	initRoot()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initRoot prepares the root command, sets viper defaults.
// 2nd part of initialization - initApp() - is called when cobra is initialized.
func initRoot() {
	initLogger(zerolog.InfoLevel)

	viper.SetEnvPrefix(EnvVarPrefix)
	viper.AutomaticEnv()

	cobra.OnInitialize(initApp)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "turn on verbose mode")

	// Commands
	rootCmd.AddCommand(normalize.Cmd())
	rootCmd.AddCommand(convert.Cmd())
	rootCmd.AddCommand(fetch.Cmd())
}

// initApp is called by cobra as soon as it is initialized.
// It reads the config, prepares required directories
func initApp() {
	if debug {
		initLogger(zerolog.DebugLevel)
	}
}

// initLogger configures the logger for the console
func initLogger(logLevel zerolog.Level) {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, FormatTimestamp: func(any) string { return "" }}).
		Level(logLevel)
}
