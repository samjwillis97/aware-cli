// Package init contains the command for generating a configuration.
package init

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"ampaware.com/cli/internal/config"
	"ampaware.com/cli/internal/utils"
)

type initParams struct {
	server       string
	organisation string
	login        string
	password     string
	authProvider string
	force        bool
}

// NewCmdInit is an init command.
func NewCmdInit() *cobra.Command {
	cmd := cobra.Command{
		Use:     "init",
		Short:   "Init initializes aware config",
		Long:    "Init initializes aware configuration required for the tool to work properly.",
		Aliases: []string{"initialize", "configure", "config", "setup"},
		Run:     initialize,
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().String("server", "", "Link to the aware api")
	cmd.Flags().String("login", "", "Aware login username or email")
	cmd.Flags().String("password", "", "Aware login password")
	cmd.Flags().String("organisation", "", "Your default organisation id")
	cmd.Flags().String("provider", "", "Authentication provider to use")
	cmd.Flags().Bool("force", false, "Forcefully override existing config if it exists")

	return &cmd
}

func getFlags(flags *pflag.FlagSet) *initParams {
	server, err := flags.GetString("server")
	utils.ExitIfError(err)

	login, err := flags.GetString("login")
	utils.ExitIfError(err)

	password, err := flags.GetString("password")
	utils.ExitIfError(err)

	organisation, err := flags.GetString("organisation")
	utils.ExitIfError(err)

	provider, err := flags.GetString("provider")
	utils.ExitIfError(err)

	force, err := flags.GetBool("force")
	utils.ExitIfError(err)

	return &initParams{
		server:       server,
		login:        login,
		password:     password,
		organisation: organisation,
		authProvider: provider,
		force:        force,
	}
}

func initialize(cmd *cobra.Command, _ []string) {
	params := getFlags(cmd.Flags())

	c := config.NewAwareCLIConfigGenerator(
		&config.AwareCLIConfig{
			Server:       params.server,
			Organisation: params.organisation,
			Login:        params.login,
			Password:     params.password,
			AuthProvider: params.authProvider,
			Force:        params.force,
		},
	)

	file, err := c.Generate()
	if err != nil {
		fmt.Println()
		utils.Failed("Unable to generate configuration: %s", err.Error())
		os.Exit(1)
	}

	utils.Success("Configuration generated: %s", file)
}
