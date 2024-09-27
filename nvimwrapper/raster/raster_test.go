package raster

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	raster := New()
	raster.Resize(20, 1)
	raster.CursorGoto(0, 0)
	raster.Put(0, 0, []rune("                    "))
	raster.Put(0, 0, []rune("fuck bauer"))
	lines := raster.Render()

	text := strings.Join(lines, "\n")
	require.Equal(t, "fuck bauer          ", text)
}
