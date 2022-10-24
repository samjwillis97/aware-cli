package view

import (
	"fmt"

	"ampaware.com/cli/internal/utils"
	iview "ampaware.com/cli/internal/view"
	"ampaware.com/cli/pkg/aware"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: Consider interfacing
type viewParams struct {
	ID    string
	plain bool
}

// TODO: Consider interfacing
type viewCmd struct {
	client  *aware.Client
	params  *viewParams
	device  *aware.Device
	devices []*aware.Device
}

func NewCmdView() *cobra.Command {
	cmd := cobra.Command{
		Use:     "view ID",
		Short:   "View a device",
		Long:    "See Above",        // TODO: Fix
		Example: "Should make some", // TODO: Fix
		Aliases: []string{"open"},
		Run:     view,
	}

	return &cmd
}

func SetFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("plain", false, "Display output in plain mode")
}

func view(cmd *cobra.Command, args []string) {
	server := viper.GetString("server")
	token := viper.GetString("token")

	params := parseFlagsAndArgs(cmd, args)

	client := aware.NewClient(aware.Config{
		Server:   server,
		Token:    token,
		Insecure: true,
		Debug:    viper.GetBool("debug"),
	})

	view := viewCmd{
		client: client,
		params: params,
	}

	if view.params.ID == "" {
		utils.ExitIfError(view.setDevices())
		utils.ExitIfError(view.getDeviceID())
	} else {
		utils.ExitIfError(view.setDevice())
	}

	// TODO: Replace this with the view
	v := iview.DeviceView{
		Data: view.device,
		Display: iview.DeviceViewDisplayFormat{
			Plain:            false,
			ExcludedSections: make(map[string]struct{}),
		},
	}

	v.Render()
}

func (v *viewCmd) setDevices() error {
	s := utils.ShowLoading("Fetching Devices...")
	defer s.Stop()

	devices, err := v.client.GetAllDevices(aware.GetAllDevicesOptions{
		OrganisationID: viper.GetString("organisation"),
	})
	if err != nil {
		return err
	}

	v.devices = devices
	return nil
}

func (v *viewCmd) getDeviceID() error {
	var ans string

	options := make([]string, 0)
	for _, device := range v.devices {
		options = append(options, device.ID+" - "+device.ParentEntity.GetParentHierachyName()+" - "+device.DisplayName)
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

	for _, device := range v.devices {
		if ans == device.ID+" - "+device.ParentEntity.GetParentHierachyName()+" - "+device.DisplayName {
			v.params.ID = device.ID
			v.device = device
			break
		}
	}

	return nil
}

func (v *viewCmd) setDevice() error {
	s := utils.ShowLoading(fmt.Sprintf("Fetching Device %s", v.params.ID))
	defer s.Stop()

	device, err := v.client.GetDeviceByID(v.params.ID)
	if err != nil {
		return err
	}

	v.device = device

	return nil
}

func parseFlagsAndArgs(cmd *cobra.Command, args []string) *viewParams {
	plain, err := cmd.Flags().GetBool("plain")
	utils.ExitIfError(err)

	var id string
	if len(args) >= 1 {
		id = args[0]
	}

	return &viewParams{
		ID:    id,
		plain: plain,
	}
}
