package raster

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderSingleLine(t *testing.T) {
	raster := New[rune]()
	raster.Resize(20, 1)
	raster.fillWith(' ')
	raster.Put(0, 0, []rune("Hello"))
	lines := RenderStringArray(raster)
	require.Equal(t, []string{"Hello               "}, lines)
}

func TestRenderMultipleLines(t *testing.T) {
	raster := New[rune]()
	raster.Resize(15, 3)
	raster.fillWith(' ')
	raster.Put(0, 0, []rune("Line1"))
	raster.Put(1, 0, []rune("Line2"))
	raster.Put(2, 0, []rune("Line3"))
	lines := RenderStringArray(raster)
	expected := []string{
		"Line1          ",
		"Line2          ",
		"Line3          ",
	}
	require.Equal(t, expected, lines)
}

func TestOverwriteLine(t *testing.T) {
	raster := New[rune]()
	raster.Resize(10, 1)
	raster.fillWith(' ')
	raster.Put(0, 0, []rune("Test"))
	raster.Put(0, 0, []rune("Over"))
	lines := RenderStringArray(raster)
	require.Equal(t, []string{"Over      "}, lines)
}

func TestPartialOverwrite(t *testing.T) {
	raster := New[rune]()
	raster.Resize(10, 1)
	raster.fillWith(' ')
	raster.Put(0, 0, []rune("12345"))
	raster.Put(0, 2, []rune("abc"))
	lines := RenderStringArray(raster)
	require.Equal(t, []string{"12abc     "}, lines)
}

func TestFullLineOverwrite(t *testing.T) {
	raster := New[rune]()
	raster.Resize(5, 1)
	raster.fillWith(' ')
	raster.Put(0, 0, []rune("Hello"))
	lines := RenderStringArray(raster)
	require.Equal(t, []string{"Hello"}, lines)
}

func TestEmptyLine(t *testing.T) {
	raster := New[rune]()
	raster.Resize(5, 1)
	raster.fillWith(' ')
	lines := RenderStringArray(raster)
	require.Equal(t, []string{"     "}, lines)
}

func TestWriteAtDifferentPositions(t *testing.T) {
	raster := New[rune]()
	raster.Resize(20, 3)
	raster.fillWith(' ')
	raster.Put(0, 0, []rune("Start"))
	raster.Put(1, 5, []rune("Middle"))
	raster.Put(2, 15, []rune("End"))
	lines := RenderStringArray(raster)
	expected := []string{
		"Start               ",
		"     Middle         ",
		"               End  ",
	}
	require.Equal(t, expected, lines)
}

func TestMoveRegionUp(t *testing.T) {
	raster := New[rune]()
	raster.Resize(4, 5)
	raster.Put(0, 0, []rune("fkbr"))
	raster.Put(1, 0, []rune("FKBR"))
	raster.Put(2, 0, []rune("RBKF"))
	raster.Put(3, 0, []rune("SXOE"))
	raster.Put(4, 0, []rune("sxoe"))

	raster.ScrollRegion(BoundingBox{
		Top:   0,
		Bot:   4,
		Left:  0,
		Right: 5,
	}, 1)

	lines := RenderStringArray(raster)

	require.Equal(t, []string{
		"FKBR",
		"RBKF",
		"SXOE",
		"SXOE",
		"sxoe"}, lines)
}

func TestMoveRegionDown(t *testing.T) {
	raster := New[rune]()
	raster.Resize(4, 5)
	raster.Put(0, 0, []rune("fkbr"))
	raster.Put(1, 0, []rune("FKBR"))
	raster.Put(2, 0, []rune("RBKF"))
	raster.Put(3, 0, []rune("SXOE"))
	raster.Put(4, 0, []rune("sxoe"))

	raster.ScrollRegion(BoundingBox{
		Top:   0,
		Bot:   4,
		Left:  0,
		Right: 5,
	}, -1)

	lines := RenderStringArray(raster)

	require.Equal(t, []string{
		"fkbr",
		"fkbr",
		"FKBR",
		"RBKF",
		"sxoe"}, lines)
}
