// Package root contains the setup and root of the aware CLI.
package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"ampaware.com/cli/internal/cmd/device"
	initCmd "ampaware.com/cli/internal/cmd/init"
	awareConfig "ampaware.com/cli/internal/config"
	"ampaware.com/cli/internal/utils"
)

var (
	config string
	debug  bool
)

func init() {
	cobra.OnInitialize(func() {
		if config != "" {
			// Use the config supplied via argument if it is supplied
			viper.SetConfigFile(config)
		} else {
			// Read the config from the default directory
			configDir, err := awareConfig.GetConfigDirectory()
			if err != nil {
				utils.Failed("Error Finding configuration directory: %v", err)
				return
			}

			// Sets up the config directory, filename, and filetype
			viper.AddConfigPath(configDir)
			viper.SetConfigName(awareConfig.ConfigFileName)
			viper.SetConfigType(awareConfig.ConfigFileType)
		}

		// Load any environment keys that match configured
		viper.AutomaticEnv()

		// Defines a prefix for environment variables to use
		// - helps to avoid clashing with other programs
		viper.SetEnvPrefix("aware")

		// Load the config file from disk
		if err := viper.ReadInConfig(); err == nil && debug {
			fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
		}
	})
}

// NewCmdRoot is the root command of the aware CLI.
func NewCmdRoot() *cobra.Command {
	cmd := cobra.Command{
		Use:   "aware <command> <subcommand>",
		Short: "Interactive AWARE CLI.",
		Long:  "Interactive AWARE Command Line Interface.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// This will be execute when a command is run but returns an error
			return cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Children of this command will inherit and execute this
			subCmd := cmd.Name()
			if !cmdRequireToken(subCmd) {
				// If a command doesn't require a token skip checking
				return
			}

			configFile := viper.ConfigFileUsed()
			if !awareConfig.Exists(configFile) {
				utils.Failed("Missing configuration file.\nRun 'aware init' to configure the tool.")
			}

			awareConfig.CheckForToken()
		},
	}

	configDir, err := awareConfig.GetConfigDirectory()
	if err != nil {
		utils.Failed("Error getting config dir: %v", err)
	}

	// Persistent flags are available to every child command of this command
	cmd.PersistentFlags().StringVarP(
		&config, "config", "c", "",
		fmt.Sprintf("Config file (default is %s%s%s.yml)", configDir, string(os.PathSeparator), awareConfig.ConfigFileName),
	)
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Turn on debug output")

	// This allows the overwriting of viper config with the flag given to cobra
	_ = viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))

	addChildCommands(&cmd)

	return &cmd
}

func addChildCommands(cmd *cobra.Command) {
	// TODO: Switch Organisation
	// TODO: Automatic re-auth
	cmd.AddCommand(
		initCmd.NewCmdInit(),
		device.NewCmdDevice(),
	)
}

func cmdRequireToken(cmd string) bool {
	allowList := []string{
		"aware",
		"init",
		"help",
		"version",
		"completion",
		"man",
	}

	for _, item := range allowList {
		if item == cmd {
			return false
		}
	}

	return true
}
