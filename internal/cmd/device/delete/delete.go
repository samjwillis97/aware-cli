// Package delete contains the command for deleting a device.
package delete

import (
	"fmt"

	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/pkg/aware"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type deleteParams struct {
	ID    string
	force bool
}

type deleteCommand struct {
	client  *aware.Client
	params  *deleteParams
	devices []*aware.Device
}

// NewCmdDelete is the delete device command.
func NewCmdDelete() *cobra.Command {
	return &cobra.Command{
		Use:     "delete ID",
		Short:   "Delete a device",
		Long:    "See Above",        // TODO: Fix
		Example: "Should make some", // TODO: Fix
		Aliases: []string{"remove", "rm", "del"},
		Run:     del,
	}
}

func del(cmd *cobra.Command, args []string) {
	server := viper.GetString("server")
	token := viper.GetString("token")

	params := parseFlagsAndArgs(cmd, args)

	client := aware.NewClient(aware.Config{
		Server:   server,
		Token:    token,
		Insecure: true,
		Debug:    viper.GetBool("debug"),
	})

	del := deleteCommand{
		client: client,
		params: params,
	}

	if del.params.ID == "" {
		utils.ExitIfError(del.setDevices())
		utils.ExitIfError(del.getDeviceID())
	}

	if !del.params.force {
		var confirm bool

		qs := &survey.Question{
			Name:     "id",
			Prompt:   &survey.Confirm{Message: fmt.Sprintf("Are you sure you want to delete %s?", del.params.ID)},
			Validate: survey.Required,
		}

		if err := survey.Ask([]*survey.Question{qs}, &confirm); err != nil {
			utils.ExitIfError(err)
		}

		if !confirm {
			return
		}
	}

	err := func() error {
		s := utils.ShowLoading(fmt.Sprintf("Removing Device %s", del.params.ID))
		defer s.Stop()

		err := del.client.DeleteDevice(del.params.ID)
		if err != nil {
			return err
		}

		return nil
	}()

	utils.ExitIfError(err)

	utils.Success("Device removed successfully\n%s", del.params.ID)
}

// SetFlags set the flags supported by the the delete command.
func SetFlags(cmd *cobra.Command) {
	// TODO: Could do a cascade when doing this for entities
	cmd.Flags().BoolP("force", "f", false, "Force the deletion of the device.")
}

func parseFlagsAndArgs(cmd *cobra.Command, args []string) *deleteParams {
	force, err := cmd.Flags().GetBool("force")
	utils.ExitIfError(err)

	var id string
	if len(args) >= 1 {
		id = args[0]
	}

	return &deleteParams{
		ID:    id,
		force: force,
	}
}

func (d *deleteCommand) setDevices() error {
	devices, err := d.client.GetAllDevices(aware.GetAllDevicesOptions{})
	if err != nil {
		return err
	}

	d.devices = devices
	return nil
}

func (d *deleteCommand) getDeviceID() error {
	var ans string

	options := make([]string, 0)
	for _, device := range d.devices {
		options = append(options, device.ID+" - "+device.DisplayName)
	}

	qs := &survey.Question{
		Name: "id",
		Prompt: &survey.Select{
			Message: "Device:",
			Options: options,
		},
		Validate: survey.Required,
	}

	if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
		return err
	}

	for _, device := range d.devices {
		if ans == device.ID+" - "+device.DisplayName {
			d.params.ID = device.ID
			break
		}
	}

	return nil
}
