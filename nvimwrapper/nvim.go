package nvimwrapper

import (
	"fmt"
	"log/slog"

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
	err = v.AttachUI(Cols, Rows, make(map[string]interface{}))
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
		case "resize":
			va := update[1].([]interface{})
			ia := make([]int, 0, 2)
			for _, v := range va {
				ia = append(ia, int(v.(int64)))
			}
			n.r.Resize(ia[0], ia[1])
		case "cursor_goto":
			va := update[1].([]interface{})
			ia := make([]int, 0, 2)
			for _, v := range va {
				ia = append(ia, int(v.(int64)))
			}
			n.r.CursorGoto(ia[0], ia[1])
		case "put":
			a := make([]rune, 0, len(update)-1)
			for _, v := range update[1:] {
				ia := v.([]interface{})
				for _, rv := range ia {
					s := rv.(string)
					a = append(a, []rune(s)[0])
				}
			}
			n.r.Put(a)
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
