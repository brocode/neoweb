package raster

import "log/slog"

type Raster struct {
	raster [][]rune
	Row    int
	Col    int
}

func New() *Raster {
	return &Raster{}
}

func (r *Raster) Resize(cols, rows int) {
	slog.Info("Resize raster", "rows", rows, "cols", cols)
	r.raster = make([][]rune, rows)
	for i := range r.raster {
		r.raster[i] = make([]rune, cols)
	}
}

func (r *Raster) CursorGoto(row, col int) {
	slog.Info("Cursor Goto", "row", row, "col", col)
	r.Row = row
	r.Col = col
}

func (r *Raster) Put(rowIdx, colIdx int, runes []rune) {
	slog.Info("Put", "text", string(runes))
	row := r.raster[rowIdx]
	copy(row[colIdx:], runes)
}

func (r *Raster) Render() []string {
	lines := make([]string, 0, len(r.raster))

	for _, row := range r.raster {
		line := string(row)
		lines = append(lines, line)
	}

	return lines
}
