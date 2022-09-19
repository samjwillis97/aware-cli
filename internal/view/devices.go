package view

import (
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/pkg/aware"
	"ampaware.com/cli/pkg/tui/table"
	tea "github.com/charmbracelet/bubbletea"
)

type DeviceDisplayFormat struct {
	Plain      bool
	NoHeaders  bool
	Columns    []string
	NoTruncate bool
}

type DeviceList struct {
	Total   int
	Server  string
	Data    []*aware.Device
	Display DeviceDisplayFormat
	Refresh func() ([]*aware.Device, error)
	// Refresh Function for TUI
	// See the pkgs/tuis
}

func (d *DeviceList) Render() error {
	if d.Display.Plain {
		w := tabwriter.NewWriter(os.Stdout, 0, tabWidth, 1, '\t', 0)
		return d.renderPlain(w)
	}

	cols, rows := d.getTableFormattedData()

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithAutoWidth(true),
		table.WithFullscreen(true),
		table.WithRefresh(d.getTableFormattedData),
		table.WithFocused(true))

	p := tea.NewProgram(t)

	if err := p.Start(); err != nil {
		utils.Failed("Error has occurred: %v", err)
	}
	return nil
}

func (d *DeviceList) getTableFormattedData() ([]table.Column, []table.Row) {
	devices, err := d.Refresh()
	utils.ExitIfError(err)
	d.Data = devices
	data := d.data()

	var cols []table.Column
	var rows []table.Row
	for i, row := range data {
		if i == 0 {
			for _, col := range row {
				cols = append(cols, table.Column{Title: col, Width: 10})
			}
		} else {
			rows = append(rows, row)
		}
	}

	return cols, rows
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
		headers   []string
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

// TODO: maybe change to tui.TableData.
func (d *DeviceList) data() [][]string {
	data := make([][]string, 0)

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
	return []string{
		fieldUID,
		fieldDisplayName,
		fieldType,
		fieldDescription,
		fieldParent,
		fieldEnabled,
	}
}
