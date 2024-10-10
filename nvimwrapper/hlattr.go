package nvimwrapper

import (
	"fmt"
	"strings"
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
func (hl HlAttr) TextDecoration() string {
	decorations := []string{}
	if hl.Underline {
		decorations = append(decorations, "underline")
	}

	if hl.Strikethrough {
		decorations = append(decorations, "line-through")
	}
	return strings.Join(decorations, " ")
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
