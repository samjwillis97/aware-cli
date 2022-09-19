package view

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/charmbracelet/glamour"
)

const (
	wordWrap = 120
	tabWidth = 8
	helpText = ""
)

// TODO: maybe change to tui.TableData.
func renderPlain(w io.Writer, data [][]string) error {
	for _, items := range data {
		n := len(items)
		for j, v := range items {
			fmt.Fprintf(w, "%s", v)
			if j != n-1 {
				fmt.Fprintf(w, "\t")
			}
		}
		fmt.Fprintln(w)
	}

	if _, ok := w.(*tabwriter.Writer); ok {
		return w.(*tabwriter.Writer).Flush()
	}
	return nil
}

func MDRenderer() (*glamour.TermRenderer, error) {
	return glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
		glamour.WithWordWrap(wordWrap),
	)
}
