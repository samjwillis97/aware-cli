package generate

import (
	"fmt"
	"time"

	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/pkg/aware"
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

    printHeaders(device.DeviceType.Parameters)
    ts, publishedValues := publishParameterValues(client, device)
    printValues(ts, publishedValues)

    if (singleValue) {return}

    timeTicker := (
        time.Duration(frequencySeconds) * time.Second + 
        time.Duration(frequencyMinutes) * time.Minute)

    ticker := time.NewTicker(timeTicker)
    quit := make(chan struct{})
    for {
        select {
        case <- ticker.C:
            ts, publishedValues := publishParameterValues(client, device)
            printValues(ts, publishedValues)
        case <- quit:
            ticker.Stop()
            return
        }
    }
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

func SetFlags(cmd *cobra.Command) {
    cmd.Flags().BoolP("single-value", "s", false, "Only generates a single value for each parameter")
    cmd.Flags().Int("frequency-seconds", 30, "The second frequency in which to generate values")
    cmd.Flags().Int("frequency-minutes", 0, "The minute frequency in which to generate values")
}
