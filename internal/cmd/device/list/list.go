// Package list contains the command for listing all devices.
package list

import (
	"fmt"

	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/internal/view"
	"ampaware.com/cli/pkg/aware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCmdList is the command for listing devices.
func NewCmdList() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List lists devices in an organisation",
		Long:    "See Above",        // TODO: Fix
		Example: "Should make some", // TODO: Fix
		Aliases: []string{"lists", "ls"},
		Run:     List,
	}
}

// List is Run by NewCmdList to load the list of devices.
func List(cmd *cobra.Command, _ []string) {
	loadList(cmd)
}

func loadList(cmd *cobra.Command) {
	devices, total, err := func() ([]*aware.Device, int, error) {
		s := utils.ShowLoading("Fetching Devices...")
		defer s.Stop()
		resp, err := loadDevices()
		return resp, len(resp), err
	}()
	utils.ExitIfError(err)

	if total == 0 {
		fmt.Println()
		utils.Failed("No results found for given query")
		return
	}

	// Handle Flags like Plain here
	plain, err := cmd.Flags().GetBool("plain")
	utils.ExitIfError(err)

	noHeaders, err := cmd.Flags().GetBool("no-headers")
	utils.ExitIfError(err)

	noTruncate, err := cmd.Flags().GetBool("no-truncate")
	utils.ExitIfError(err)

	// receiving := make(chan table.TransmitMessage)
	// transmit := make(chan table.ReceiveMessage)

	v := view.DeviceList{
		Total:  total,
		Server: viper.GetString("server"),
		Data:   devices,
		Display: view.DeviceDisplayFormat{
			Plain:      plain,
			NoHeaders:  noHeaders,
			NoTruncate: noTruncate,
		},
		Refresh: loadDevices,
		// Transmit: receiving,
		// Receive: transmit,
	}

	// Maybe move this to the view layer
	// Just feed in the required functions?
	// Could use the With.. Format or just as straight functions

	// Yeah this needs to be done on the view, to get the correct columns etc.
	// go func() {
	//     for {
	//         switch <-receiving {
	//         case table.RefreshPressed:
	//             devices, err := loadDevices()

	//             utils.ExitIfError(err)
	//             v.Data = devices
	//             transmit <- table.RefreshCompleted
	//         case table.Exit:
	//             break
	//         }
	//     }
	// }()

	utils.ExitIfError(v.Render())
}

func loadDevices() ([]*aware.Device, error) {
	client := aware.NewClient(aware.Config{
		Server:   viper.GetString("server"),
		Token:    viper.GetString("token"),
		Insecure: true,
		Debug:    viper.GetBool("debug"),
	})

	resp, err := client.GetAllDevices(aware.GetAllDevicesOptions{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetFlags sets all the flags for the command.
func SetFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("plain", false, "Display output in plain mode")
	cmd.Flags().Bool("no-truncate", false, "Show all available columns in plain mode. Works only with --plain")
	cmd.Flags().Bool("no-headers", false, "Don't display headers in plain mode. Works only with --plain")
}
