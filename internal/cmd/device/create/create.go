// Package create contains the command for creating a new device.
package create

import (
	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/pkg/aware"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type createParams struct {
	deviceType   string
	parentEntity string
	organisation string
	displayName  string
	noInput      bool
}

type createCmd struct {
	client         *aware.Client
	params         *createParams
	deviceTypes    []*aware.DeviceType
	parentEntities []*aware.Entity
}

// NewCmdCreate is the create device command.
func NewCmdCreate() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Short:   "Create a new device",
		Long:    "See Above",        // TODO: Fix
		Example: "Should make some", // TODO: Fix
		Run:     create,
	}
}

// SetFlags sets the flags support by the create command.
func SetFlags(cmd *cobra.Command) {
	cmd.Flags().String("type", "", "Set the Device Type.")
	cmd.Flags().String("parent", "", "Set the Parent Entity.")
	cmd.Flags().String("name", "", "Set the Display Name.")
}

func create(cmd *cobra.Command, _ []string) {
	server := viper.GetString("server")
	token := viper.GetString("token")

	params := parseFlags(cmd)

	client := aware.NewClient(aware.Config{
		Server:   server,
		Token:    token,
		Insecure: true,
		Debug:    viper.GetBool("debug"),
	})

	cc := createCmd{
		client: client,
		params: params,
	}

	if cc.isNonInteractive() {
		// TODO: Handle
		cc.params.noInput = true
	}

	utils.ExitIfError(cc.setDeviceTypes())
	utils.ExitIfError(cc.setParentEntities())
	utils.ExitIfError(cc.askQuestions())

	ID, err := func() (string, error) {
		s := utils.ShowLoading("Creating an issue...")
		defer s.Stop()

		cr := aware.CreateDeviceRequest{
			DisplayName:  params.displayName,
			Organisation: params.organisation,
			ParentEntity: params.parentEntity,
			DeviceType:   params.deviceType,
		}

		resp, err := client.CreateDevice(&cr)
		if err != nil {
			return "", err
		}
		return resp.ID, err
	}()
	utils.ExitIfError(err)
	utils.Success("Device created\n%s", ID)
}

func (c *createCmd) setDeviceTypes() error {
	deviceTypes, err := c.client.GetAllDeviceTypes(c.params.organisation)
	if err != nil {
		return err
	}

	c.deviceTypes = deviceTypes
	return nil
}

func (c *createCmd) setParentEntities() error {
	parentEntities, err := c.client.GetAllEntities(c.params.organisation, aware.GetAllEntitiesOptions{})
	if err != nil {
		return err
	}

	c.parentEntities = parentEntities
	return nil
}

func (c *createCmd) getDeviceType() *survey.Question {
	var qs *survey.Question

	if c.params.deviceType != "" {
		return qs
	}

	options := make([]string, 0)
	for _, t := range c.deviceTypes {
		options = append(options, t.Name)
	}

	qs = &survey.Question{
		Name: "deviceType",
		Prompt: &survey.Select{
			Message: "Device type:",
			Options: options,
		},
		Validate: survey.Required,
	}

	return qs
}

func (c *createCmd) getParentEntity() *survey.Question {
	var qs *survey.Question

	if c.params.parentEntity != "" {
		return qs
	}

	options := make([]string, 0)
	options = append(options, "(None)")
	for _, t := range c.parentEntities {
		options = append(options, t.GetParentHierachyName())
	}

	qs = &survey.Question{
		Name: "parentEntity",
		Prompt: &survey.Select{
			Message: "Parent entity:",
			Options: options,
		},
		Validate: survey.Required,
	}

	return qs
}

func (c *createCmd) askQuestions() error {
	deviceType := c.getDeviceType()
	if deviceType != nil {
		ans := struct{ DeviceType string }{}
		err := survey.Ask([]*survey.Question{deviceType}, &ans)
		if err != nil {
			return err
		}

		if c.params.deviceType == "" {
			for _, t := range c.deviceTypes {
				if t.Name == ans.DeviceType {
					c.params.deviceType = t.ID
					break
				}
			}
		}
	}

	parentEntity := c.getParentEntity()
	if parentEntity != nil {
		ans := struct{ ParentEntity string }{}
		err := survey.Ask([]*survey.Question{parentEntity}, &ans)
		if err != nil {
			return err
		}

		if c.params.parentEntity == "" {
			for _, t := range c.parentEntities {
				if t.GetParentHierachyName() == ans.ParentEntity {
					c.params.parentEntity = t.ID
					break
				}
			}
		}
	}

	var qs []*survey.Question

	if c.params.displayName == "" {
		qs = append(qs, &survey.Question{
			Name:   "displayName",
			Prompt: &survey.Input{Message: "Display Name:"},
		})
	}

	ans := struct{ DisplayName string }{}
	err := survey.Ask(qs, &ans)
	if err != nil {
		return err
	}

	if c.params.displayName == "" {
		c.params.displayName = ans.DisplayName
	}

	return nil
}

func (c *createCmd) isNonInteractive() bool {
	return utils.StdinHasData()
}

func (c *createCmd) isMandatoryParmMissing() bool {
	return (c.params.organisation == "" ||
		c.params.deviceType == "" ||
		c.params.parentEntity == "")
}

func parseFlags(cmd *cobra.Command) *createParams {
	deviceType, err := cmd.Flags().GetString("type")
	utils.ExitIfError(err)

	parentEntity, err := cmd.Flags().GetString("parent")
	utils.ExitIfError(err)

	displayName, err := cmd.Flags().GetString("name")
	utils.ExitIfError(err)

	return &createParams{
		deviceType:   deviceType,
		parentEntity: parentEntity,
		displayName:  displayName,
		organisation: viper.GetString("organisation"),
	}
}
