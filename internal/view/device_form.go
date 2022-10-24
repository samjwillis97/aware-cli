package view

import (
	"fmt"

	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/pkg/aware"
	"ampaware.com/cli/pkg/tui/form"
	tea "github.com/charmbracelet/bubbletea"
)

type DeviceViewDisplayFormat struct {
	Plain            bool
	ExcludedSections map[string]struct{}
}

type DeviceView struct {
	Data     *aware.Device
	sections []form.Section
	Display  DeviceViewDisplayFormat
}

func (d *DeviceView) Render() error {
	if d.sections == nil {
		d.sections = d.getSections()
	}

	if d.Display.Plain {
		fmt.Println("Please Implement Me")
		return nil
	}

	f := form.New(
		form.WithSections(d.sections),
	)

	p := tea.NewProgram(f)

	if err := p.Start(); err != nil {
		utils.Failed("Error has occurred: %v", err)
	}

	return nil
}

func (d *DeviceView) getSections() []form.Section {
	sections := make([]form.Section, 0)

	if _, have := d.Display.ExcludedSections["Info"]; !have {
		infoSection := form.Section{
			Name: "Info",
			Fields: []*form.Field{
				{
					Name:  "ID",
					Value: d.Data.ID,
				},
				{
					Name:  "Name",
					Value: d.Data.DisplayName,
				},
				{
					Name:  "Device Type",
					Value: d.Data.DeviceType.Name,
				},
				{
					Name:  "Parent Entity",
					Value: d.Data.ParentEntity.GetParentHierachyName(),
				},
				{
					Name:  "Cloud ID",
					Value: d.Data.CloudID,
				},
			},
		}

		sections = append(sections, infoSection)
	}

	return sections
}
