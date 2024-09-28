package raster

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderSingleLine(t *testing.T) {
	raster := New()
	raster.Resize(20, 1)
	raster.fillWithSpaces()
	raster.Put(0, 0, []rune("Hello"))
	lines := raster.Render()
	require.Equal(t, []string{"Hello               "}, lines)
}

func TestRenderMultipleLines(t *testing.T) {
	raster := New()
	raster.Resize(15, 3)
	raster.fillWithSpaces()
	raster.Put(0, 0, []rune("Line1"))
	raster.Put(1, 0, []rune("Line2"))
	raster.Put(2, 0, []rune("Line3"))
	lines := raster.Render()
	expected := []string{
		"Line1          ",
		"Line2          ",
		"Line3          ",
	}
	require.Equal(t, expected, lines)
}

func TestOverwriteLine(t *testing.T) {
	raster := New()
	raster.Resize(10, 1)
	raster.fillWithSpaces()
	raster.Put(0, 0, []rune("Test"))
	raster.Put(0, 0, []rune("Over"))
	lines := raster.Render()
	require.Equal(t, []string{"Over      "}, lines)
}

func TestPartialOverwrite(t *testing.T) {
	raster := New()
	raster.Resize(10, 1)
	raster.fillWithSpaces()
	raster.Put(0, 0, []rune("12345"))
	raster.Put(0, 2, []rune("abc"))
	lines := raster.Render()
	require.Equal(t, []string{"12abc     "}, lines)
}

func TestFullLineOverwrite(t *testing.T) {
	raster := New()
	raster.Resize(5, 1)
	raster.fillWithSpaces()
	raster.Put(0, 0, []rune("Hello"))
	lines := raster.Render()
	require.Equal(t, []string{"Hello"}, lines)
}

func TestEmptyLine(t *testing.T) {
	raster := New()
	raster.Resize(5, 1)
	raster.fillWithSpaces()
	lines := raster.Render()
	require.Equal(t, []string{"     "}, lines)
}

func TestWriteAtDifferentPositions(t *testing.T) {
	raster := New()
	raster.Resize(20, 3)
	raster.fillWithSpaces()
	raster.Put(0, 0, []rune("Start"))
	raster.Put(1, 5, []rune("Middle"))
	raster.Put(2, 15, []rune("End"))
	lines := raster.Render()
	expected := []string{
		"Start               ",
		"     Middle         ",
		"               End  ",
	}
	require.Equal(t, expected, lines)
}

