package generate

import (
	"fmt"
	"os"
	"time"

	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/internal/view"
	"ampaware.com/cli/pkg/aware"
	"ampaware.com/cli/pkg/tui/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdDeviceTelemetryGenerate() *cobra.Command {
    cmd := cobra.Command{
        Use: "generate",
        Short: "Generate telemetry for a device",
        Long: "Generate telemetry for a Device", // TODO: Fix
        Example: "Should do something here", // TODO: Fix
        Aliases: []string{},
        Annotations: map[string]string{},
        Args: cobra.MinimumNArgs(1),
        Run: generate,
        // TODO: Help for Args
    }

    return &cmd
}

func generate(cmd *cobra.Command, args []string) {
    // Default will generate a telemetry value for each parameter once per 10 seconds
    deviceID := args[0]

    singleValue, err := cmd.Flags().GetBool("single-value")
    utils.ExitIfError(err)

    frequencySeconds, err := cmd.Flags().GetInt("frequency-seconds")
    utils.ExitIfError(err)

    frequencyMinutes, err := cmd.Flags().GetInt("frequency-minutes")
    utils.ExitIfError(err)

    client := aware.NewClient(aware.Config{
        Server:   viper.GetString("server"),
        Token: viper.GetString("token"),
        Insecure: true,
        Debug:    viper.GetBool("debug"),
    })

    device, err := func() (*aware.Device, error) {
        s := utils.ShowLoading("Fetching Device...")
        defer s.Stop()
    
        resp, err := client.GetDeviceByID(deviceID); if err != nil {
            return nil, err
        } else {
            return resp, nil
        }
    }()
    utils.ExitIfError(err)

    appendReady := make(chan byte)
    appendRow, err := publishValuesToRow(client, device)
    utils.ExitIfError(err)

    t := view.TelemetryTable{
        Parameters: &device.DeviceType.Parameters,
        Display: view.TelemetryTableDisplayFormat{
            Plain: false,
            NoHeaders: false,
        },
        AppendRow: &appendRow,
        AppendReady: appendReady,
        InitialRows: []table.Row{appendRow},
    }

    if (!singleValue) {
        timeTicker := (
            time.Duration(frequencySeconds) * time.Second + 
            time.Duration(frequencyMinutes) * time.Minute)
        ticker := time.NewTicker(timeTicker)
        signalChan := make(chan os.Signal, 1)
        quit := make(chan struct{})
        go func() {
            for {
                select {
                case <- ticker.C:
                    appendRow, err = publishValuesToRow(client, device)
                    utils.ExitIfError(err)
                    appendReady<- 1
                case <- signalChan:
                    ticker.Stop()
                    return
                case <- quit:
                    ticker.Stop()
                    return
                }
            }
        }()
    }

    utils.ExitIfError(t.Render())
}

func printHeaders(parameters []aware.DeviceTypeParameter) {
    toPrint := "Generated Values for "
    for _, val := range parameters {
        toPrint += fmt.Sprintf("%s    ", val.DisplayName)
    }
    fmt.Print(toPrint+"\n")
}

func printValues(ts time.Time, values []interface{}) {
    toPrint := ts.Format(time.RFC3339) + "    "
    for _, val := range values {
        toPrint += fmt.Sprintf("%v    ", val)
    }
    fmt.Print(toPrint+"\n")
}

func publishParameterValues(client *aware.Client, device *aware.Device) (time.Time, []interface{}) {
    ts := time.Now()
    var publishedValues []interface{}
    for _, parameter := range device.DeviceType.Parameters {
        value := parameter.GetRandomValue()
        publishedValues = append(publishedValues, value)
        utils.ExitIfError(client.PublishTelemetry(
            device.ID,
            parameter.Name,
            value,
            ts,
        ))
    }
    return ts, publishedValues
}

func publishValuesToRow(client *aware.Client, device *aware.Device) (table.Row, error) {
    ts, values := publishParameterValues(client, device)

    var row table.Row
    row = append(row, ts.Format(time.RFC3339))
    for _, val := range values {
        row = append(row, fmt.Sprintf("%v", val))
    }

    return row, nil
}

func SetFlags(cmd *cobra.Command) {
    cmd.Flags().BoolP("single-value", "s", false, "Only generates a single value for each parameter")
    cmd.Flags().Int("frequency-seconds", 30, "The second frequency in which to generate values")
    cmd.Flags().Int("frequency-minutes", 0, "The minute frequency in which to generate values")
}
