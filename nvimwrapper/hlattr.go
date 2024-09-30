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
