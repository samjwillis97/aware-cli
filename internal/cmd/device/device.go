package device

import (
	"ampaware.com/cli/internal/cmd/device/list"
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

    lc := list.NewCmdList()

    cmd.AddCommand(
        lc,
    )

    list.SetFlags(lc)

    return &cmd
}

func device(cmd *cobra.Command, _ []string) error {
    return cmd.Help()
}
