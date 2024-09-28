package nvimwrapper

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"github.com/brocode/neoweb/key"
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
		nvim.ChildProcessArgs("--embed", "--clean", "--cmd", "set noswapfile"),
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

	err = v.RegisterHandler("redraw", wrapper.handleRedraw)
	if err != nil {
		return nil, fmt.Errorf("Failed to register handler: %w", err)

	}

	return &wrapper, nil
}

func (n *NvimWrapper) handleRedraw(events ...[]interface{}) {
	for _, event := range events {
		eventName, ok := event[0].(string)
		if !ok {
			continue
		}
		updates := event[1:]

		slog.Debug("Redraw Event", "name", eventName, "updates", updates)
		for _, update := range updates {
			switch eventName {
			case "grid_resize":
				n.handleResize(update.([]interface{}))
			case "grid_cursor_goto":
				n.handleGoto(update.([]interface{}))
			case "grid_line":
				n.handleGridLine(update.([]interface{}))
			}
		}
	}
}

func (n *NvimWrapper) handleGridLine(line_data []interface{}) {
	// TODO grid id is ignored for now
	row := line_data[1].(int64)
	col := line_data[2].(int64)

	slog.Debug("put grid_line", "line", line_data)
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
	slog.Debug("put grid_line interpreted", "row", row, "col", col, "text", buffer.String())
	n.r.Put(int(row), int(col), []rune(buffer.String()))

}
func (n *NvimWrapper) handleGoto(update []interface{}) {
	ia := make([]int, 0, 2)
	// TODO first element is the grid id, multiple grids are supported
	for _, v := range update[1:] {
		ia = append(ia, int(v.(int64)))
	}
	n.r.CursorGoto(ia[0], ia[1])
}
func (n *NvimWrapper) handleResize(update []interface{}) {
	ia := make([]int, 0, 2)
	// TODO first element is the grid id, multiple grids are supported
	for _, v := range update[1:] {
		ia = append(ia, int(v.(int64)))
	}
	n.r.Resize(ia[0], ia[1])
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

func (w *NvimWrapper) SendKey(keyPress key.KeyPress) {
	// TODO actually consider modifiers. "shift" is also a key
	err := w.Input(keyPress.Key)

	if err != nil {
		slog.Error("Failed to process keypress.", "press", keyPress)
	}

}

func (w *NvimWrapper) Render() (NvimResult, error) {

	lines := w.r.Render()

	return NvimResult{
		Lines:          lines,
		CursorPosition: [2]int{w.r.Row, w.r.Col},
	}, nil
}
