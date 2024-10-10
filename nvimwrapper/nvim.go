package nvimwrapper

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/brocode/neoweb/key"
	"github.com/brocode/neoweb/nvimwrapper/raster"
	"github.com/neovim/go-client/nvim"
)

const (
	Rows = 40
	Cols = 120
)

type hlRune struct {
	rune rune
	hlId int
}

type Span struct {
	Text string
	HlId int
}
type RenderedLine struct {
	Spans []Span
}

type NvimResult struct {
	Lines          []RenderedLine
	CursorPosition [2]int
}

type NvimWrapper struct {
	v    *nvim.Nvim
	r    *raster.Raster[hlRune]
	hl   map[int]HlAttr
	cond *sync.Cond
	mu   sync.Mutex
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

func Spawn(clean bool) (*NvimWrapper, error) {

	args := []string{"--embed", "--cmd", "set noswapfile"}
	if clean {
		args = append(args, "--clean")
	}
	// Start an embedded Neovim process
	v, err := nvim.NewChildProcess(
		nvim.ChildProcessArgs(args...),
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

	wrapper := NvimWrapper{
		r:  raster.New[hlRune](),
		hl: make(map[int]HlAttr),
		v:  v,
	}
	wrapper.cond = sync.NewCond(&wrapper.mu)

	err = v.RegisterHandler("redraw", wrapper.handleRedraw)
	if err != nil {
		return nil, fmt.Errorf("Failed to register handler: %w", err)

	}

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
	slog.Debug("Send input", "input", input)
	_, err := w.v.Input(input)
	if err != nil {
		return fmt.Errorf("Failed to input: %w", err)
	}
	return nil
}

func modifiers(input string, pressed key.KeyPress) string {
	modifiers := []string{}
	if pressed.CtrlKey {
		modifiers = append(modifiers, "C")
	}
	if pressed.AltKey {
		modifiers = append(modifiers, "M")
	}
	// TODO dunno if this will fuck us
	if pressed.ShiftKey && len(input) > 1 && pressed.Key != "<" && pressed.Key != ">" {
		modifiers = append(modifiers, "S")
	}

	if len(modifiers) < 1 {
		return input
	}
	modifierStr := strings.Join(modifiers, "-")
	if len(input) > 1 {
		return fmt.Sprintf("<%s-%s>", modifierStr, input[1:len(input)-1])
	} else {
		return fmt.Sprintf("<%s-%s>", modifierStr, input)
	}
}

func (w *NvimWrapper) SendKey(keyPress key.KeyPress) {
	input := ""
	// TODO handle modifiers
	switch keyPress.Key {
	case "Escape":
		input = "<Esc>"
	case "Enter":
		input = "<CR>"
	case "Tab":
		input = "<Tab>"
	case "Backspace":
		input = "<BS>"
	case "Delete":
		input = "<Del>"
	case "ArrowUp":
		input = "<Up>"
	case "ArrowDown":
		input = "<Down>"
	case "ArrowLeft":
		input = "<Left>"
	case "ArrowRight":
		input = "<Right>"
	case "Home":
		input = "<Home>"
	case "End":
		input = "<End>"
	case "PageUp":
		input = "<PageUp>"
	case "PageDown":
		input = "<PageDown>"
	case "Insert":
		input = "<Insert>"
	case "<":
		input = "<LT>"
	default:
		input = keyPress.Key
	}
	err := w.Input(modifiers(input, keyPress))

	if err != nil {
		slog.Error("Failed to process keypress.", "press", keyPress)
	}
}

func (n *NvimWrapper) Render() (NvimResult, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.render()
}

func (n *NvimWrapper) render() (NvimResult, error) {
	lines := renderHlRunes(n.r)

	return NvimResult{
		Lines:          lines,
		CursorPosition: [2]int{n.r.Row, n.r.Col},
	}, nil
}

func renderHlRunes(r *raster.Raster[hlRune]) []RenderedLine {
	lines := []RenderedLine{}

	for _, row := range r.Raster {
		line := RenderedLine{}

		hlId := row[0].hlId
		spanBuffer := []rune{}

		for _, currentRune := range row {
			if currentRune.hlId != hlId {
				line.Spans = append(line.Spans, Span{
					Text: string(spanBuffer),
					HlId: hlId,
				})
				hlId = currentRune.hlId
				spanBuffer = []rune{}
			}
			spanBuffer = append(spanBuffer, currentRune.rune)

		}
		if len(spanBuffer) > 0 {
			line.Spans = append(line.Spans, Span{
				Text: string(spanBuffer),
				HlId: hlId,
			})
		}

		lines = append(lines, line)
	}

	return lines
}

func (n *NvimWrapper) RenderOnFlush(ctx context.Context, handler func(result NvimResult) error) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// wake up in case the client disconnected
	go func() {
		<-ctx.Done()
		n.mu.Lock()
		defer n.mu.Unlock()

		n.cond.Broadcast()
	}()

	for {
		result, err := n.render()
		if err != nil {
			return err
		}
		err = handler(result)
		if err != nil {
			return err
		}

		n.cond.Wait()

		select {
		case <-ctx.Done():
			return nil
		default:
			continue
		}
	}

}
