package hl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	hl := &HlAttr{}
	hl.Apply("background", int64(0xFFFFFF))
	hl.Apply("foreground", int64(0x000000))
	hl.Apply("bold", true)
	hl.Apply("underline", true)
	hl.Apply("reverse", true)
	hl.Apply("italic", true)
	hl.Apply("strikethrough", true)
	hl.Apply("blend", int64(50))
	hl.Apply("special", int64(0xFF00FF))
	hl.Apply("undercurl", true)

	assert.Equal(t, "#FF00FF", *hl.Special, "Expected Special to be set to #ff00ff")
	assert.NotNil(t, hl.Blend, "Expected Blend to be non-nil")
	assert.Equal(t, 50, *hl.Blend, "Expected Blend to be set to 50")
	assert.True(t, hl.Strikethrough, "Expected Strikethrough to be true")
	assert.True(t, hl.Italic, "Expected Italic to be true")
	assert.True(t, hl.Underline, "Expected Underline to be true")
	assert.True(t, hl.Reverse, "Expected Reverse to be true")
	assert.True(t, hl.Bold, "Expected Bold to be true")
	assert.Equal(t, "#000000", *hl.Foreground, "Expected Foreground to be set to #000000")
	assert.Equal(t, "#FFFFFF", *hl.Background, "Expected Background to be set to #ffffff")
	assert.True(t, hl.Undercurl, "Expected Undercurl to be true")
}
