package telemetry

import (
	"ampaware.com/cli/internal/cmd/device/telemetry/generate"
	"github.com/spf13/cobra"
)

func NewCmdDeviceTelemetry() *cobra.Command {
	cmd := cobra.Command{
		Use:         "telemetry",
		Short:       "Manage Telemetry for a Device",
		Long:        "Manage Telemetry for a Device", // TODO: Fix
		Aliases:     []string{},
		Annotations: map[string]string{},
		RunE:        telemetry,
	}

	gen := generate.NewCmdDeviceTelemetryGenerate()

	cmd.AddCommand(
		gen,
	)

	generate.SetFlags(gen)

	return &cmd
}

func telemetry(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
