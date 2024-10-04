package nvimwrapper

import (
	"testing"

	"github.com/brocode/neoweb/key"
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
