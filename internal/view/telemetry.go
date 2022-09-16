package view

import (
	"ampaware.com/cli/pkg/aware"
)

type TelemetryTableDisplayFormat struct {
    Columns [] string
}

type TelemetryTable struct {
    Parameters []*aware.DeviceTypeParameter
    Display TelemetryTableDisplayFormat
}

func (t *TelemetryTable) Render() error {
    return nil
}
