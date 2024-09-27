package nvimwrapper

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"github.com/brocode/neoweb/nvimwrapper/raster"
	"github.com/neovim/go-client/nvim"
)

const (
	Rows = 40
	Cols = 120
)

type NvimResult struct {
	Lines          []string
	CursorPosition [2]int
}

type NvimWrapper struct {
	v *nvim.Nvim
	r *raster.Raster
}

type Line struct {
	Text   string
	Number int
}

func (r NvimResult) Row() int {
	return r.CursorPosition[0]
}

func (r NvimResult) Col() int {
	return r.CursorPosition[1]
}

func Spawn() (*NvimWrapper, error) {

	wrapper := NvimWrapper{}
	wrapper.r = raster.New()

	// Start an embedded Neovim process
	v, err := nvim.NewChildProcess(
		nvim.ChildProcessArgs("--embed", "--clean"),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to start embedded neovim: %w", err)
	}

	// Set UI dimensions (rows and columns)
	attachConfig := map[string]interface{}{"ext_linegrid": true}
	err = v.AttachUI(Cols, Rows, attachConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to attach UI: %w", err)

	}
	wrapper.v = v

	v.RegisterHandler("redraw", wrapper.handleRedraw)

	return &wrapper, nil
}

func (n *NvimWrapper) handleRedraw(updates ...[]interface{}) {
	for _, update := range updates {
		eventName, ok := update[0].(string)
		if !ok {
			continue
		}

		slog.Info("Redraw Event", "name", eventName)
		switch eventName {
		case "grid_resize":
			va := update[1].([]interface{})
			ia := make([]int, 0, 2)
			// TODO first element is the grid id, multiple grids are supported
			for _, v := range va[1:] {
				ia = append(ia, int(v.(int64)))
			}
			n.r.Resize(ia[0], ia[1])
		case "grid_cursor_goto":
			va := update[1].([]interface{})
			ia := make([]int, 0, 2)
			// TODO first element is the grid id, multiple grids are supported
			for _, v := range va[1:] {
				ia = append(ia, int(v.(int64)))
			}
			n.r.CursorGoto(ia[0], ia[1])
		case "grid_line":
			for _, line := range update[1:] {
				//["grid_line", grid, row, col_start, cells, wrap]
				//Redraw a continuous part of a row on a grid, starting at the column col_start.
				line_data := line.([]interface{})
				// TODO grid id is ignored for now
				row := line_data[1].(int64)
				col := line_data[2].(int64)

				slog.Info("put grid_line", "line", line)
				var buffer bytes.Buffer
				// cells is an array of arrays each with 1 to 3 items: [text(, hl_id, repeat)]
				for _, cell := range line_data[3].([]interface{}) {
					cell_contents := cell.([]interface{})
					text := cell_contents[0].(string)
					if len(cell_contents) == 3 {
						text = strings.Repeat(text, int(cell_contents[2].(int64)))
					}
					buffer.WriteString(text)
				}
				slog.Info("put grid_line interpreted", "row", row, "col", col, "text", buffer.String())
				// TODO should not move cursor for this
				n.r.CursorGoto(int(row), int(col))
				n.r.Put([]rune(buffer.String()))
			}
		}
	}
}

func (w *NvimWrapper) Close() {
	err := w.v.DetachUI()
	if err != nil {
		slog.Error("Failed to detach UI", "Error", err)
	}
	err = w.v.Close()
	if err != nil {
		slog.Error("Failed to close neovim", "Error", err)
	}
}

func (w *NvimWrapper) OpenFile(file string) error {
	err := w.v.Command(fmt.Sprintf("edit %v", file))
	if err != nil {
		return fmt.Errorf("Failed to open file in neovim: %w", err)
	}
	return nil
}

func (w *NvimWrapper) Input(input string) error {
	_, err := w.v.Input(input)
	if err != nil {
		return fmt.Errorf("Failed to input: %w", err)
	}
	return nil
}

func (w *NvimWrapper) Render() (NvimResult, error) {

	lines := w.r.Render()

	return NvimResult{
		Lines:          lines,
		CursorPosition: [2]int{w.r.Row, w.r.Col},
	}, nil
}
