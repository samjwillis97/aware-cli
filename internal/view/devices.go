package view

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"ampaware.com/cli/pkg/aware"
)

type DeviceDisplayFormat struct {
    Plain bool
    NoHeaders bool
    Columns [] string
    NoTruncate bool
}

type DeviceList struct {
    Total int
    Server string
    Data []*aware.Device
    Display DeviceDisplayFormat
    FooterText string
    // Refresh Function for TUI
    // See the pkgs/tuis
}


func (d *DeviceList) Render() error {
    if d.Display.Plain {
        w := tabwriter.NewWriter(os.Stdout, 0, tabWidth, 1, '\t', 0)
        return d.renderPlain(w)
    }

    data := d.data()
    if d.FooterText == "" {
        d.FooterText = fmt.Sprintf("Showing %d of %d results for devices", len(data)-1, d.Total)
    }

    return nil
}

func (d *DeviceList) renderPlain(w io.Writer) error {
    return renderPlain(w, d.data())
}

func (d *DeviceList) header() []string {
    if len(d.Display.Columns) == 0 {
        validColumns := ValidDeviceColumns()
        if d.Display.NoTruncate || !d.Display.Plain {
            return validColumns
        }
        if len(validColumns) > 4 {
            return validColumns[0:4] // Why 0-4????, is this just for nicely displaying
        }
        return validColumns
    }

    var (
        headers []string
        hasUIDCol bool
    )

    columnsSet := d.validColumnsSet()
    for _, c := range d.Display.Columns {
        c = strings.ToUpper(c)
        if _, ok := columnsSet[c]; ok {
            headers = append(headers, strings.ToUpper(c))
        }
        if c == fieldUID {
            hasUIDCol = true
        }
    }

    if !hasUIDCol {
        headers = append([]string{fieldUID}, headers...)
    }

    return headers
}

// TODO: maybe change to tui.TableData
func (d *DeviceList) data() [][]string {
    var data [][]string

    headers := d.header()

    if !(d.Display.Plain && d.Display.NoHeaders) {
        data = append(data, headers)
    }

    if len(headers) == 0 {
        headers = ValidDeviceColumns()
    }

    for _, device := range d.Data {
        data = append(data, d.assignColumns(headers, device))
    }

    return data
}

func (d *DeviceList) validColumnsSet() map[string]struct{} {
    columns := ValidDeviceColumns()
    out := make(map[string]struct{}, len(columns))

    for _, c := range columns {
        out[c] = struct{}{}
    }

    return out
}

func (DeviceList) assignColumns(columns []string, device *aware.Device) []string {
    var bucket []string

    for _, column := range columns {
        switch column {
        case fieldUID:
            bucket = append(bucket, device.ID)
        case fieldDisplayName:
            bucket = append(bucket, device.DisplayName)
        case fieldType:
            bucket = append(bucket, device.DeviceType.Name)
        case fieldDescription:
            bucket = append(bucket, device.DeviceType.Description)
        case fieldParent:
            bucket = append(bucket, device.ParentEntity.GetParentHierachyName())
        case fieldEnabled:
            bucket = append(bucket, strconv.FormatBool(device.IsEnabled))
        }
    }

    return bucket
}

func ValidDeviceColumns() []string {
    // TODO: Finish This
    return []string{
        fieldUID,
        fieldDisplayName,
        fieldType,
        fieldDescription,
        fieldParent,
        fieldEnabled,
    }
}
