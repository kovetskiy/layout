package main

import (
	"strings"

	"github.com/hpcloud/golor"
)

const (
	KeyHeight = 3
)

const (
	DefaultKeyWidth = 3

	CornerTopLeft     = "┌"
	LineHorizontal    = "─"
	CornerTopRight    = "┐"
	LineVertical      = "│"
	CornerBottomRight = "┘"
	CornerBottomLeft  = "└"
)

type Key struct {
	Name   string
	Label  string
	Width  int
	Margin int
	Hold   bool
	Mod    bool

	X int
	Y int
}

func (key *Key) Render(style *Style) [KeyHeight]string {
	lines := [KeyHeight]string{}

	width := key.Width
	if width == 0 {
		width = DefaultKeyWidth
	}

	label := key.Label
	if label == "" {
		label = key.Name
	}

	labelLen := len(label)
	if width < len(label) {
		width = labelLen
	}

	var paddingLeft string
	var paddingRight string
	empty := width - labelLen
	if empty > 0 {
		paddingLeft = strings.Repeat(" ", empty/2)
		// we have empty - empty / 2 here in case when empty % 2 == 1
		// it will mean that on left side we will have one space less than on
		// right side
		paddingRight = strings.Repeat(" ", empty-empty/2)
	}

	margin := strings.Repeat(" ", key.Margin)

	lines[0] = CornerTopLeft +
		strings.Repeat(LineHorizontal, width) +
		CornerTopRight

	lines[1] = LineVertical +
		paddingLeft +
		label +
		paddingRight +
		LineVertical

	lines[2] = CornerBottomLeft +
		strings.Repeat(LineHorizontal, width) +
		CornerBottomRight

	for i := 0; i < KeyHeight; i++ {
		line := lines[i]
		if style != nil && key.Hold {
			line = colorize(line, style, key.Mod)
		}

		lines[i] = margin + line
	}

	return lines
}

func colorize(line string, style *Style, mod bool) string {
	if mod {
		return golor.Colorize(
			line,
			style.Mod.Hold.Foreground,
			style.Mod.Hold.Background,
		)
	}

	return golor.Colorize(
		line,
		style.Tap.Hold.Foreground,
		style.Tap.Hold.Background,
	)
}

func (key *Key) GetWidth() int {
	width := key.Width
	if width == 0 {
		width = DefaultKeyWidth
	}

	label := key.Label
	if label == "" {
		label = key.Name
	}

	labelLen := len(label)
	if width < len(label) {
		width = labelLen
	}

	width += 2 // vertical lines / corners
	width += key.Margin

	return width
}
