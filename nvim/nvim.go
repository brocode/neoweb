package nvim

import (
	"strings"

	"github.com/neovim/go-client/nvim"
)

type NvimResult struct {
	VisibleText    string
	CursorPosition [2]int
}

func RunNvim() (NvimResult, error) {
	// Start an embedded Neovim process
	v, err := nvim.NewChildProcess(nvim.ChildProcessArgs("--embed"))
	if err != nil {
		return NvimResult{}, err
	}
	defer v.Close()

	// Set UI dimensions (rows and columns)
	err = v.AttachUI(80, 24, make(map[string]interface{}))
	if err != nil {
		return NvimResult{}, err

	}
	defer v.DetachUI()

	// Open a file
	err = v.Command("edit main.go")
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

	// Get cursor position
	pos, err := v.WindowCursor(win)
	if err != nil {
		return NvimResult{}, err

	}

	stringLines := make([]string, 0, len(lines))
	for _, line := range lines {
		stringLines = append(stringLines, string(line))
	}

	return NvimResult{
		VisibleText:    strings.Join(stringLines, "\n"),
		CursorPosition: pos,
	}, nil

}
