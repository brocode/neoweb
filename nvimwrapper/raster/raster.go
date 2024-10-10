package raster

import "log/slog"

type BoundingBox struct {
	Top   int
	Bot   int
	Left  int
	Right int
}

type Raster[T any] struct {
	Raster [][]T
	Row    int
	Col    int
}

func New[T any]() *Raster[T] {
	return &Raster[T]{}
}

func (r *Raster[T]) Resize(cols, rows int) {
	slog.Debug("Resize raster", "rows", rows, "cols", cols)
	r.Raster = make([][]T, rows)
	for i := range r.Raster {
		r.Raster[i] = make([]T, cols)
	}
}

func (r *Raster[T]) fillWith(fill T) {
	for i := range r.Raster {
		for j := range r.Raster[i] {
			r.Raster[i][j] = fill
		}
	}
}

func (r *Raster[T]) CursorGoto(row, col int) {
	slog.Debug("Cursor Goto", "row", row, "col", col)
	r.Row = row
	r.Col = col
}

func (r *Raster[T]) Put(rowIdx, colIdx int, runes []T) {
	row := r.Raster[rowIdx]
	copy(row[colIdx:], runes)
}

func RenderStringArray(r *Raster[rune]) []string {
	lines := make([]string, 0, len(r.Raster))

	for _, row := range r.Raster {
		line := string(row)
		lines = append(lines, line)
	}

	return lines
}

func (r *Raster[T]) ScrollRegion(boundingBox BoundingBox, rowMovement int) {
	if rowMovement > 0 {
		for rowIdx := boundingBox.Top + rowMovement; rowIdx < boundingBox.Bot; rowIdx++ {
			sliceToMove := r.Raster[rowIdx][boundingBox.Left : boundingBox.Right-1]
			destinationSlice := r.Raster[rowIdx-rowMovement][boundingBox.Left : boundingBox.Right-1]
			copy(destinationSlice, sliceToMove)
		}
	} else {
		for rowIdx := boundingBox.Bot + rowMovement - 1; rowIdx >= boundingBox.Top; rowIdx-- {
			sliceToMove := r.Raster[rowIdx][boundingBox.Left : boundingBox.Right-1]
			destinationSlice := r.Raster[rowIdx-rowMovement][boundingBox.Left : boundingBox.Right-1]
			copy(destinationSlice, sliceToMove)
		}
	}
}
