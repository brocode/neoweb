package raster

import "log/slog"

type BoundingBox struct {
	Top   int
	Bot   int
	Left  int
	Right int
}

type Raster[T any] struct {
	raster [][]T
	Row    int
	Col    int
}

func New[T any]() *Raster[T] {
	return &Raster[T]{}
}

func (r *Raster[T]) Resize(cols, rows int) {
	slog.Debug("Resize raster", "rows", rows, "cols", cols)
	r.raster = make([][]T, rows)
	for i := range r.raster {
		r.raster[i] = make([]T, cols)
	}
}

func (r *Raster[T]) fillWith(fill T) {
	for i := range r.raster {
		for j := range r.raster[i] {
			r.raster[i][j] = fill
		}
	}
}

func (r *Raster[T]) CursorGoto(row, col int) {
	slog.Debug("Cursor Goto", "row", row, "col", col)
	r.Row = row
	r.Col = col
}

func (r *Raster[T]) Put(rowIdx, colIdx int, runes []T) {
	row := r.raster[rowIdx]
	copy(row[colIdx:], runes)
}

func RenderStringArray(r *Raster[rune]) []string {
	lines := make([]string, 0, len(r.raster))

	for _, row := range r.raster {
		line := string(row)
		lines = append(lines, line)
	}

	return lines
}

func (r *Raster[T]) ScrollRegion(boundingBox BoundingBox, rowMovement int) {
	if rowMovement > 0 {
		for rowIdx := boundingBox.Top + rowMovement; rowIdx < boundingBox.Bot; rowIdx++ {
			sliceToMove := r.raster[rowIdx][boundingBox.Left : boundingBox.Right-1]
			destinationSlice := r.raster[rowIdx-rowMovement][boundingBox.Left : boundingBox.Right-1]
			copy(destinationSlice, sliceToMove)
		}
	} else {
		for rowIdx := boundingBox.Bot + rowMovement - 1; rowIdx >= boundingBox.Top; rowIdx-- {
			sliceToMove := r.raster[rowIdx][boundingBox.Left : boundingBox.Right-1]
			destinationSlice := r.raster[rowIdx-rowMovement][boundingBox.Left : boundingBox.Right-1]
			copy(destinationSlice, sliceToMove)
		}
	}
}
