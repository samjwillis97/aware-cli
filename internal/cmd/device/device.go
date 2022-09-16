package device

import (
	"ampaware.com/cli/internal/cmd/device/list"
	"ampaware.com/cli/internal/cmd/device/telemetry"
	"github.com/spf13/cobra"
)

func NewCmdDevice() *cobra.Command {
    cmd := cobra.Command{
        Use: "device",
        Short: "Manage Devices in an Organisation",
        Long: "See Above?", // TODO: Fix
        Aliases: []string{"devices"},
        Annotations: map[string]string{"": ""}, // TODO: What is this?
        RunE: device,
    }

    // TODO: Create
    // TODO: Delete -> Are You Sure?
    // TODO: Edit
    // TODO: View
    // TODO: State?
    // TODO: Telemetry ->
    //  - Watch
    //  - Generate
    // TODO: Parameter, Telemetry
    //  - Watch
    //  - Generate

    lc := list.NewCmdList()

    cmd.AddCommand(
        lc,
        telemetry.NewCmdDeviceTelemetry(),
    )

    list.SetFlags(lc)

    return &cmd
}

func device(cmd *cobra.Command, _ []string) error {
    return cmd.Help()
}
