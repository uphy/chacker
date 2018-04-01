package result

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type (
	GridResultBody struct {
		header []string
		rows   [][]string
	}
)

func NewGridResultBody(header ...string) *GridResultBody {
	return &GridResultBody{header, [][]string{}}
}

func (g *GridResultBody) Append(row ...string) {
	g.rows = append(g.rows, row)
}

func (g *GridResultBody) JSON() interface{} {
	rows := []map[string]string{}
	for _, row := range g.rows {
		r := map[string]string{}
		for i, cell := range row {
			name := g.header[i]
			r[name] = cell
		}
		rows = append(rows, r)
	}
	return rows
}

func (g *GridResultBody) Pretty(writer io.Writer) error {
	w := tablewriter.NewWriter(writer)
	w.SetHeader(g.header)
	for _, row := range g.rows {
		w.Append(row)
	}
	w.Render()
	return nil
}

func (g *GridResultBody) Plain(writer io.Writer) error {
	fmt.Fprintln(writer, strings.Join(g.header, " "))
	for _, row := range g.rows {
		fmt.Fprintln(writer, strings.Join(row, " "))
	}
	return nil
}
