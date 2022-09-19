package view

import (
	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/pkg/aware"
	"ampaware.com/cli/pkg/tui/table"
	tea "github.com/charmbracelet/bubbletea"
)

type TelemetryTableDisplayFormat struct {
    Plain bool
    NoHeaders bool
}

type TelemetryTable struct {
    Parameters *[]aware.DeviceTypeParameter
    Display TelemetryTableDisplayFormat
    AppendRow *table.Row
    AppendReady chan byte
    InitialRows []table.Row
}

func (v *TelemetryTable) Render() error {
    if v.Display.Plain {
        return nil
    }

    cols := v.getColumns()

    t := table.New(
        table.WithColumns(cols),
        table.WithRows(v.InitialRows),
        table.WithAutoWidth(true),
        table.WithFullscreen(true),
        table.WithFocused(true),
        table.WithAppending(v.AppendRow),
    )

    p := tea.NewProgram(t)

    go func() {
        for {
            <- v.AppendReady
            p.Send(table.AppendReady)
        }
    }()

    if err := p.Start(); err != nil {
        utils.Failed("Error has occurred: %v", err)
    }

    return nil
}

func (v *TelemetryTable) getColumns() []table.Column {
    var cols []table.Column
    cols = append(cols, table.Column{"Time", 10})
    for _, val := range *v.Parameters {
        cols = append(cols, table.Column{val.DisplayName, 10})
    }

    return cols
}
