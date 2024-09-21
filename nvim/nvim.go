package nvim

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

const (
	Rows = 40
	Cols = 120
)

type NvimResult struct {
	Lines          []Line
	CursorPosition [2]int
	Status         string
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

func RunNvim() (NvimResult, error) {
	// Start an embedded Neovim process
	v, err := nvim.NewChildProcess(
		nvim.ChildProcessArgs("--embed", "--clean"),
	)
	if err != nil {
		return NvimResult{}, err
	}
	defer v.Close()

	// Set UI dimensions (rows and columns)
	err = v.AttachUI(Cols, Rows, make(map[string]interface{}))
	if err != nil {
		return NvimResult{}, err

	}
	defer v.DetachUI()

	// Open a file
	err = v.Command("edit main.go")
	if err != nil {
		return NvimResult{}, err

	}

	// Move cursor.
	// just for now. this needs to come from the frontend
	_, err = v.Input("ggjj$")
	if err != nil {
		return NvimResult{}, err

	}

	// Get the current buffer and window
	buf, err := v.CurrentBuffer()
	if err != nil {
		return NvimResult{}, err

	}
	win, err := v.CurrentWindow()
	if err != nil {
		return NvimResult{}, err

	}

	// Get the first and last visible lines
	var firstLine int
	err = v.Eval("line('w0')", &firstLine)
	if err != nil {
		return NvimResult{}, err

	}
	var lastLine int
	err = v.Eval("line('w$')", &lastLine)
	if err != nil {
		return NvimResult{}, err

	}

	// Adjust for 0-based indexing in Go
	firstLineNum := firstLine - 1 // Vimscript lines are 1-based
	lastLineNum := lastLine       // No need to adjust end index

	// Get visible lines
	lines, err := v.BufferLines(buf, firstLineNum, lastLineNum, true)
	if err != nil {
		return NvimResult{}, err

	}

	stringLines := make([]Line, 0, len(lines))
	for idx, line := range lines {
		stringLines = append(stringLines, Line{
			Text:   string(line),
			Number: idx + firstLine,
		})
	}

	// Get cursor position
	pos, err := v.WindowCursor(win)
	if err != nil {
		return NvimResult{}, err
	}
	status, err := getStatus(v)
	if err != nil {
		return NvimResult{}, err
	}

	return NvimResult{
		Lines:          stringLines,
		CursorPosition: pos,
		Status:         status,
	}, nil

}

func getStatus(v *nvim.Nvim) (string, error) {
	// Evaluate the statusline with default options
	result, err := v.EvalStatusLine("%{mode()} %f %h%m%r%=%-14.(%l,%c%V%) %P", map[string]interface{}{})
	if err != nil {
		return "", err
	}

	// Extract the "str" field from the result map
	str, ok := result["str"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get statusline string")
	}

	return strings.TrimSpace(str), nil
}
