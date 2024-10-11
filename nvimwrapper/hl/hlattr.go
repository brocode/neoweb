package hl

import (
	"fmt"
	"strings"

	"github.com/brocode/neoweb/nvimwrapper/vimnumbers"
)

type HlAttr struct {
	Background    *string
	Foreground    *string
	Blend         *int
	Special       *string
	Bold          bool
	Underline     bool
	Reverse       bool
	Italic        bool
	Strikethrough bool
	Undercurl     bool
}

func (hl HlAttr) Color() string {
	if hl.Foreground != nil {
		return *hl.Foreground
	}
	return "inherit"
}
func (hl HlAttr) BackgroundColor() string {
	if hl.Background != nil {
		return *hl.Background
	}
	return "inherit"
}
func (hl HlAttr) FontWeight() string {
	if hl.Bold {
		return "bold"
	}
	return "normal"
}
func (hl HlAttr) FontStyle() string {
	if hl.Italic {
		return "italic"
	}
	return "normal"
}
func (hl HlAttr) TextDecorationColor() string {
	if hl.Special != nil {
		return *hl.Special
	}

	return "currentcolor"
}
func (hl HlAttr) TextDecorationStyle() string {
	if hl.Undercurl {
		return "wavy"
	}

	return "solid"
}
func (hl HlAttr) TextDecoration() string {
	decorations := []string{}
	if hl.Underline || hl.Undercurl {
		decorations = append(decorations, "underline")
	}

	if hl.Strikethrough {
		decorations = append(decorations, "line-through")
	}

	if len(decorations) > 0 {
		return strings.Join(decorations, " ")
	} else {
		return "inherit"
	}
}

func (hl *HlAttr) Apply(key string, value interface{}) {
	switch key {
	case "background":
		hl.Background = vimnumbers.ConvertToHexColor(vimnumbers.ForceInt(value))
	case "foreground":
		hl.Foreground = vimnumbers.ConvertToHexColor(vimnumbers.ForceInt(value))
	case "bold":
		hl.Bold = value.(bool)
	case "underline":
		hl.Underline = value.(bool)
	case "reverse":
		hl.Reverse = value.(bool)
	case "italic":
		hl.Italic = value.(bool)
	case "strikethrough":
		hl.Strikethrough = value.(bool)
	case "blend":
		intValue := vimnumbers.ForceInt(value)
		hl.Blend = &intValue
	case "special":
		hl.Special = vimnumbers.ConvertToHexColor(vimnumbers.ForceInt(value))
	case "undercurl":
		hl.Undercurl = value.(bool)

	}
}

func (h HlAttr) String() string {
	var parts []string

	if h.Background != nil {
		parts = append(parts, fmt.Sprintf("Background: %s", *h.Background))
	}
	if h.Foreground != nil {
		parts = append(parts, fmt.Sprintf("Foreground: %s", *h.Foreground))
	}
	if h.Bold {
		parts = append(parts, "Bold")
	}
	if h.Underline {
		parts = append(parts, "Underline")
	}
	if h.Reverse {
		parts = append(parts, "Reverse")
	}
	if h.Italic {
		parts = append(parts, "Italic")
	}
	if h.Strikethrough {
		parts = append(parts, "Strikethrough")
	}
	if h.Blend != nil {
		parts = append(parts, fmt.Sprintf("Blend: %d", *h.Blend))
	}
	if h.Special != nil {
		parts = append(parts, fmt.Sprintf("Special: %s", *h.Special))
	}
	if h.Undercurl {
		parts = append(parts, "Undercurl")
	}

	return strings.Join(parts, ", ")
}
