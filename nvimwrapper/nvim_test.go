package nvimwrapper

import (
	"testing"

	"github.com/brocode/neoweb/key"
	"github.com/brocode/neoweb/nvimwrapper/raster"
	"github.com/stretchr/testify/require"
)

func TestTranslateModifiers(t *testing.T) {

	require.Equal(t,
		"<C-w>",
		modifiers("w", key.KeyPress{
			CtrlKey: true,
		}))

	require.Equal(t,
		"W",
		modifiers("W", key.KeyPress{
			ShiftKey: true,
		}))

	require.Equal(t,
		"<C-W>",
		modifiers("W", key.KeyPress{
			ShiftKey: true,
			CtrlKey:  true,
		}))

	require.Equal(t,
		"<C-S-F1>",
		modifiers("<F1>", key.KeyPress{
			ShiftKey: true,
			CtrlKey:  true,
		}))

	require.Equal(t,
		"<M-Esc>",
		modifiers("<Esc>", key.KeyPress{
			AltKey: true,
		}))

	require.Equal(t,
		"<C-M-O>",
		modifiers("O", key.KeyPress{
			AltKey:   true,
			CtrlKey:  true,
			ShiftKey: true,
		}))
}

func TestRenderHlRunes(t *testing.T) {
	r := raster.New[hlRune]()
	r.Resize(4, 1)
	r.Put(0, 0, []hlRune{
		{rune: 'f', hlId: 1},
		{rune: 'k', hlId: 1},
		{rune: 'b', hlId: 1},
		{rune: 'r', hlId: 2},
	})
	lines := renderHlRunes(r)
	require.Equal(t, 1, len(lines), "Only one line")

	line := lines[0]

	require.Equal(t, 2, len(line.Spans), "2 hl ids")
	require.Equal(t, "fkb", line.Spans[0].Text)
	require.Equal(t, 1, line.Spans[0].HlId)
	require.Equal(t, 2, line.Spans[1].HlId)
}
