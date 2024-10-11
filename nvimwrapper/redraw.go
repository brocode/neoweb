package nvimwrapper

import (
	"log/slog"
	"strings"

	"github.com/brocode/neoweb/nvimwrapper/hl"
	"github.com/brocode/neoweb/nvimwrapper/raster"
	"github.com/brocode/neoweb/nvimwrapper/vimnumbers"
)

func (n *NvimWrapper) handleRedraw(events ...[]interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, event := range events {
		eventName, ok := event[0].(string)
		if !ok {
			continue
		}
		updates := event[1:]

		slog.Debug("Redraw Event", "name", eventName, "updates", updates)
		for _, update := range updates {
			tuple, ok := update.([]interface{})
			if !ok {
				slog.Error("Update is not a tuple", "update", update)
				continue
			}
			switch eventName {
			case "grid_resize":
				n.handleResize(tuple)
			case "grid_scroll":
				n.handleScroll(tuple)
			case "grid_cursor_goto":
				n.handleGoto(tuple)
			case "grid_line":
				n.handleGridLine(tuple)
			case "hl_attr_define":
				n.handleHlAttrDefine(tuple)
			case "mode_change":
				n.handleModeChange(tuple)
			case "flush":
				slog.Debug("Flush")
				n.cond.Broadcast()
			}
		}
	}
}

func (n *NvimWrapper) handleModeChange(lineData []interface{}) {
	n.mode = lineData[0].(string)
	n.modeIdx = vimnumbers.ForceInt(lineData[1])
}

func (n *NvimWrapper) handleHlAttrDefine(lineData []interface{}) {
	if len(lineData) != 4 {
		slog.Warn("Invalid hl attr define.", "data", lineData)
		return
	}
	id := vimnumbers.ForceInt(lineData[0])
	rawAttrs := lineData[1].(map[string]interface{})

	attr := hl.HlAttr{}

	for key, value := range rawAttrs {
		attr.Apply(key, value)
	}

	n.hl[id] = attr
}

func (n *NvimWrapper) handleGridLine(line_data []interface{}) {
	hlId := 0
	row := line_data[1].(int64)
	col := line_data[2].(int64)
	data := line_data[3].([]interface{})
	buffer := make([]hlRune, 0, len(data))
	// cells is an array of arrays each with 1 to 3 items: [text(, hl_id, repeat)]
	for _, cell := range data {
		cell_contents := cell.([]interface{})
		text := cell_contents[0].(string)
		if len(cell_contents) >= 2 {
			hlId = vimnumbers.ForceInt(cell_contents[1])
		}
		if len(cell_contents) == 3 {
			text = strings.Repeat(text, int(cell_contents[2].(int64)))
		}
		for _, rune := range text {
			buffer = append(buffer, hlRune{
				rune: rune,
				hlId: hlId,
			})
		}

	}
	n.r.Put(int(row), int(col), buffer)

}
func (n *NvimWrapper) handleGoto(update []interface{}) {
	ia := make([]int, 0, 2)
	for _, v := range update[1:] {
		ia = append(ia, int(v.(int64)))
	}
	n.r.CursorGoto(ia[0], ia[1])
}
func (n *NvimWrapper) handleResize(update []interface{}) {
	ia := make([]int, 0, 2)
	for _, v := range update[1:] {
		ia = append(ia, int(v.(int64)))
	}
	n.r.Resize(ia[0], ia[1])
}

func (n *NvimWrapper) handleScroll(update []interface{}) {
	slog.Debug("scroll grid", "data", update)
	boundingBox := raster.BoundingBox{
		Top:   int(update[1].(int64)),
		Bot:   int(update[2].(int64)),
		Left:  int(update[3].(int64)),
		Right: int(update[4].(int64)),
	}
	n.r.ScrollRegion(boundingBox, int(update[5].(int64)))
}
