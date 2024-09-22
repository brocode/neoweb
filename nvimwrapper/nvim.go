package nvimwrapper

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/neovim/go-client/nvim"
)

const (
	Rows = 40
	Cols = 120
)

type NvimResult struct {
	Status         string
	Lines          []Line
	CursorPosition [2]int
}

type NvimWrapper struct {
	v *nvim.Nvim
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
	return &wrapper, nil
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
	status, err := w.getStatus()
	if err != nil {
		return NvimResult{}, err
	}

	cursorPosition, err := w.getCursorPos()
	if err != nil {
		return NvimResult{}, err
	}
	lines, err := w.getVisibleLines()
	if err != nil {
		return NvimResult{}, err
	}

	return NvimResult{
		Lines:          lines,
		CursorPosition: cursorPosition,
		Status:         status,
	}, nil
}

func (w *NvimWrapper) getCursorPos() ([2]int, error) {
	v := w.v

	win, err := v.CurrentWindow()
	if err != nil {
		return [2]int{}, fmt.Errorf("Failed to get current window: %w", err)

	}

	// Get cursor position
	pos, err := v.WindowCursor(win)
	if err != nil {
		return [2]int{}, fmt.Errorf("Failed to get cursor: %w", err)
	}
	return pos, nil
}
func (w *NvimWrapper) getVisibleLines() ([]Line, error) {
	v := w.v
	buf, err := v.CurrentBuffer()
	if err != nil {
		return []Line{}, err

	}
	// Get the first and last visible lines
	var firstLine int
	err = v.Eval("line('w0')", &firstLine)
	if err != nil {
		return []Line{}, fmt.Errorf("Failed to get first line: %w", err)

	}
	var lastLine int
	err = v.Eval("line('w$')", &lastLine)
	if err != nil {
		return []Line{}, fmt.Errorf("Failed to get last line: %w", err)

	}

	// Adjust for 0-based indexing in Go
	firstLineNum := firstLine - 1 // Vimscript lines are 1-based
	lastLineNum := lastLine       // No need to adjust end index

	// Get visible lines
	lines, err := v.BufferLines(buf, firstLineNum, lastLineNum, true)
	if err != nil {
		return []Line{}, err

	}

	result := make([]Line, 0, len(lines))
	for idx, line := range lines {
		result = append(result, Line{
			Text:   string(line),
			Number: idx + firstLine,
		})
	}
	return result, nil
}

func (w *NvimWrapper) getStatus() (string, error) {
	v := w.v
	result, err := v.EvalStatusLine("%{mode()} %f %h%m%r%=%-14.(%l,%c%V%) %P", map[string]interface{}{})
	if err != nil {
		return "", fmt.Errorf("Get status: %w", err)
	}

	str, ok := result["str"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get statusline string")
	}

	return strings.TrimSpace(str), nil
}
