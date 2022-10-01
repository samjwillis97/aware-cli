// Package edit contains the command for editing an existing device.
package edit

import (
	"fmt"

	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/pkg/aware"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type editParams struct {
	ID           string
	deviceType   string
	parentEntity string
	organisation string
	displayName  string
	noInput      bool
}

type editCmd struct {
	client         *aware.Client
	params         *editParams
	device         *aware.Device
	devices        []*aware.Device
	deviceTypes    []*aware.DeviceType
	parentEntities []*aware.Entity
}

// NewCmdEdit is the edit device command.
func NewCmdEdit() *cobra.Command {
	cmd := cobra.Command{
		Use:     "edit ID",
		Short:   "Edit a device",
		Long:    "See Above",        // TODO: Fix
		Example: "Should make some", // TODO: Fix
		Aliases: []string{"update", "modify"},
		Run:     edit,
	}

	return &cmd
}

// SetFlags set the flags supported by the edit command.
func SetFlags(cmd *cobra.Command) {
	cmd.Flags().String("type", "", "Modified Device Type.")
	cmd.Flags().String("parent", "", "Modified Parent Entity.")
	cmd.Flags().String("name", "", "Modified Display Name.")
	cmd.Flags().Bool("no-input", false, "Disable prompt for non-required fields")
}

func edit(cmd *cobra.Command, args []string) {
	server := viper.GetString("server")
	token := viper.GetString("token")

	params := parseFlagsAndArgs(cmd, args)

	client := aware.NewClient(aware.Config{
		Server:   server,
		Token:    token,
		Insecure: true,
		Debug:    viper.GetBool("debug"),
	})

	edit := editCmd{
		client: client,
		params: params,
	}

	if edit.params.ID == "" {
		utils.ExitIfError(edit.setDevices())
		utils.ExitIfError(edit.getDevice())
	} else {
		utils.ExitIfError(edit.setDevice())
	}

	if edit.params.deviceType == "" {
		utils.ExitIfError(edit.setDeviceTypes())
	}
	if edit.params.parentEntity == "" {
		utils.ExitIfError(edit.setParentEntities())
	}

	edit.askQuestions()
}

func (e *editCmd) setDevice() error {
	s := utils.ShowLoading(fmt.Sprintf("Fetching Device %s", e.params.ID))
	defer s.Stop()

	device, err := e.client.GetDeviceByID(e.params.ID)
	if err != nil {
		return err
	}

	e.device = device

	return nil
}

func (e *editCmd) setDevices() error {
	s := utils.ShowLoading("Fetching Devices...")
	defer s.Stop()

	devices, err := e.client.GetAllDevices(aware.GetAllDevicesOptions{
		OrganisationID: viper.GetString("organisation"),
	})
	if err != nil {
		return err
	}

	e.devices = devices
	return nil
}

func (e *editCmd) setDeviceTypes() error {
	s := utils.ShowLoading("Fetching Device Types...")
	defer s.Stop()

	deviceTypes, err := e.client.GetAllDeviceTypes(e.params.organisation)
	if err != nil {
		return err
	}

	e.deviceTypes = deviceTypes
	return nil
}

func (e *editCmd) setParentEntities() error {
	s := utils.ShowLoading("Fetching Entities...")
	defer s.Stop()

	parentEntities, err := e.client.GetAllEntities(e.params.organisation, aware.GetAllEntitiesOptions{})
	if err != nil {
		return err
	}

	e.parentEntities = parentEntities
	return nil
}

func (e *editCmd) getDevice() error {
	var ans string

	options := make([]string, 0)
	for _, device := range e.devices {
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

	for _, device := range e.devices {
		if ans == device.ID+" - "+device.ParentEntity.GetParentHierachyName()+" - "+device.DisplayName {
			e.params.ID = device.ID
			e.device = device
			break
		}
	}

	return nil
}

func (e *editCmd) askQuestions() {
	if e.params.noInput {
		// TODO: Handle This
		return
	}

	utils.ExitIfError(e.getDeviceType())
	utils.ExitIfError(e.getParentEntity())
	utils.ExitIfError(e.getDisplayName())

	return
}

func (e *editCmd) getDeviceType() error {
	var qs *survey.Question

	if e.params.deviceType != "" {
		return nil
	}

	options := make([]string, 0)
	for _, t := range e.deviceTypes {
		options = append(options, t.Name)
	}

	qs = &survey.Question{
		Name: "deviceType",
		Prompt: &survey.Select{
			Message: fmt.Sprintf("New device type? (Currently: %s)", e.device.DeviceType.Name),
			Options: options,
			Default: e.device.DeviceType.Name,
			Help:    "Ctrl+C to skip question and leave as current",
			Description: func(value string, index int) string {
				if value == e.device.DeviceType.Name {
					return "Current"
				}
				return ""
			},
		},
	}

	ans := struct{ DeviceType string }{}
	err := survey.Ask([]*survey.Question{qs}, &ans)
	if err != nil {
		if err == terminal.InterruptErr {
			e.params.deviceType = e.device.DeviceType.ID
			utils.Success("Keeping device type: %s", e.device.DeviceType.Name)
			fmt.Println()
			return nil
		}
		return err
	}

	for _, t := range e.deviceTypes {
		if t.Name == ans.DeviceType {
			e.params.deviceType = t.ID
			break
		}
	}

	return nil
}

func (e *editCmd) getParentEntity() error {
	var qs *survey.Question

	if e.params.parentEntity != "" {
		return nil
	}

	options := make([]string, 0)
	options = append(options, "(None)")
	for _, t := range e.parentEntities {
		options = append(options, t.GetParentHierachyName())
	}

	qs = &survey.Question{
		Name: "parentEntity",
		Prompt: &survey.Select{
			// TODO: Change Message
			Message: fmt.Sprintf("New parent entity? (Currently: %s)", e.device.ParentEntity.GetParentHierachyName()),
			Options: options,
			// TODO: Fix Default
			Default: e.device.ParentEntity.GetParentHierachyName(),
			Help:    "Ctrl+C to skip question and leave as current",
			// TODO: Fix Description
			Description: func(value string, index int) string {
				if value == e.device.ParentEntity.GetParentHierachyName() {
					return "*"
				}
				return ""
			},
		},
	}

	ans := struct{ ParentEntity string }{}
	err := survey.Ask([]*survey.Question{qs}, &ans)
	if err != nil {
		if err == terminal.InterruptErr {
			e.params.parentEntity = e.device.ParentEntity.ID
			utils.Success("Keeping parent entity: %s", e.device.ParentEntity.GetParentHierachyName())
			fmt.Println()
			return nil
		}
		return err
	}

	for _, t := range e.parentEntities {
		if t.GetParentHierachyName() == ans.ParentEntity {
			e.params.parentEntity = t.ID
			break
		}
	}

	return nil
}

func (e *editCmd) getDisplayName() error {
	var qs *survey.Question

	if e.params.displayName != "" {
		return nil
	}

	qs = &survey.Question{
		Name: "displayName",
		Prompt: &survey.Input{
			Message: "New display name?",
			Default: e.device.DisplayName,
			Help:    "Ctrl+C to skip question and leave as current",
		},
	}

	ans := struct{ DisplayName string }{}
	err := survey.Ask([]*survey.Question{qs}, &ans)
	if err != nil {
		if err == terminal.InterruptErr {
			e.params.displayName = e.device.DisplayName
			utils.Success("Keeping display name: %s", e.device.DisplayName)
			fmt.Println()
			return nil
		}
		return err
	}

	e.params.displayName = ans.DisplayName

	return nil
}

func parseFlagsAndArgs(cmd *cobra.Command, args []string) *editParams {
	deviceType, err := cmd.Flags().GetString("type")
	utils.ExitIfError(err)

	parentEntity, err := cmd.Flags().GetString("parent")
	utils.ExitIfError(err)

	displayName, err := cmd.Flags().GetString("name")
	utils.ExitIfError(err)

	noInput, err := cmd.Flags().GetBool("no-input")
	utils.ExitIfError(err)

	var id string
	if len(args) >= 1 {
		id = args[0]
	}

	return &editParams{
		ID:           id,
		deviceType:   deviceType,
		parentEntity: parentEntity,
		displayName:  displayName,
		organisation: viper.GetString("organisation"),
		noInput:      noInput,
	}
}
