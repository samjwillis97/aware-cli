// Package device contains the root command for all device commands.
package device

import (
	"ampaware.com/cli/internal/cmd/device/create"
	"ampaware.com/cli/internal/cmd/device/delete"
	"ampaware.com/cli/internal/cmd/device/edit"
	"ampaware.com/cli/internal/cmd/device/list"
	"ampaware.com/cli/internal/cmd/device/telemetry"
	"github.com/spf13/cobra"
)

// NewCmdDevice is the root command for device.
func NewCmdDevice() *cobra.Command {
	cmd := cobra.Command{
		Use:         "device",
		Short:       "Manage Devices in an Organisation",
		Long:        "See Above?", // TODO: Fix
		Aliases:     []string{"devices"},
		Annotations: map[string]string{"": ""}, // TODO: What is this?
		RunE:        device,
	}

	// TODO: Register On Cloud
	// TODO: Edit
	// TODO: View
	// TODO: State?
	// TODO: Telemetry ->
	//  - Watch
	// TODO: Parameter, Telemetry
	//  - Watch
	//  - Generate

	lc := list.NewCmdList()
	cr := create.NewCmdCreate()
	de := delete.NewCmdDelete()
	ed := edit.NewCmdEdit()

	cmd.AddCommand(
		lc,
		cr,
		de,
		ed,
		telemetry.NewCmdDeviceTelemetry(),
	)

	list.SetFlags(lc)
	create.SetFlags(cr)
	delete.SetFlags(de)
	edit.SetFlags(ed)

	return &cmd
}

func device(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
