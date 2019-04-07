package main

type Row struct {
	Keys []Key
}

type Layout struct {
	Rows []Row
}

var (
	DefaultLayout = &Layout{
		Rows: []Row{
			{
				Keys: []Key{
					{Name: "ESC"},
					{Name: "F1", Margin: 5},
					{Name: "F2"},
					{Name: "F3"},
					{Name: "F4"},
					{Name: "F5", Margin: 3},
					{Name: "F6"},
					{Name: "F7"},
					{Name: "F8"},
					{Name: "F9", Margin: 3},
					{Name: "F10"},
					{Name: "F11"},
					{Name: "F12"},
					{Name: "SYSRQ", Label: "SYSRQ", Width: 6, Margin: 2},
					{Name: "SCROLLLOCK", Label: "SLOCK", Width: 6},
					{Name: "PAUSE", Label: "PAUSE", Width: 6},
				},
			},
			{
				Keys: []Key{
					{Name: "GRAVE", Label: "`"},
					{Name: "1"},
					{Name: "2"},
					{Name: "3"},
					{Name: "4"},
					{Name: "5"},
					{Name: "6"},
					{Name: "7"},
					{Name: "8"},
					{Name: "9"},
					{Name: "0"},
					{Name: "MINUS", Label: "-"},
					{Name: "EQUAL", Label: "="},
					{Name: "BACKSPACE", Width: 7},
					{Name: "INSERT", Label: "INSERT", Width: 6, Margin: 2},
					{Name: "HOME", Label: "HOME", Width: 6},
					{Name: "PAGEUP", Label: "PAGEUP", Width: 6},
				},
			},
			{
				Keys: []Key{
					{Name: "TAB", Width: 5},
					{Name: "Q"},
					{Name: "W"},
					{Name: "E"},
					{Name: "R"},
					{Name: "T"},
					{Name: "Y"},
					{Name: "U"},
					{Name: "I"},
					{Name: "O"},
					{Name: "P"},
					{Name: "LEFTBRACE", Label: "["},
					{Name: "RIGHTBRACE", Label: "]"},
					{Name: "BACKSLASH", Label: "\\", Width: 7},
					{Name: "DELETE", Label: "DELETE", Width: 6, Margin: 2},
					{Name: "END", Label: "END", Width: 6},
					{Name: "PAGEDOWN", Label: "PAGEDN", Width: 6},
				},
			},
			{
				Keys: []Key{
					{Name: "CAPSLOCK", Label: "CAPS", Width: 6, Mod: true},
					{Name: "A"},
					{Name: "S"},
					{Name: "D"},
					{Name: "F"},
					{Name: "G"},
					{Name: "H"},
					{Name: "J"},
					{Name: "K"},
					{Name: "L"},
					{Name: "SEMICOLON", Label: ";"},
					{Name: "APOSTROPHE", Label: "'"},
					{Name: "ENTER", Label: "ENTER", Width: 11},
				},
			},
			{
				Keys: []Key{
					{Name: "LEFTSHIFT", Label: "SHIFT", Width: 8, Mod: true},
					{Name: "Z"},
					{Name: "X"},
					{Name: "C"},
					{Name: "V"},
					{Name: "B"},
					{Name: "N"},
					{Name: "M"},
					{Name: "COMMA", Label: ","},
					{Name: "DOT", Label: "."},
					{Name: "SLASH", Label: "/"},
					{Name: "RIGHTSHIFT", Label: "SHIFT", Width: 14, Mod: true},
					{Name: "UP", Label: "UP", Width: 6, Margin: 10},
				},
			},
			{
				Keys: []Key{
					{Name: "LEFTCTRL", Label: "CTRL", Mod: true},
					{Name: "LEFTMETA", Label: "META", Mod: true},
					{Name: "LEFTALT", Label: "ALT", Width: 4, Mod: true},
					{Name: "SPACE", Label: "SPACE", Width: 32},
					{Name: "RIGHTALT", Label: "ALT", Width: 4, Mod: true},
					{Name: "RIGHTMETA", Label: "META", Mod: true},
					{Name: "?FN?", Label: "FN", Width: 4},
					{Name: "RIGHTCTRL", Label: "CTRL", Mod: true},
					{Name: "LEFT", Label: "LEFT", Width: 6, Margin: 2},
					{Name: "DOWN", Label: "DOWN", Width: 6},
					{Name: "RIGHT", Label: "RIGHT", Width: 6},
				},
			},
		},
	}
)

func addKeyXY(layout *Layout) {
	for indexRow, row := range layout.Rows {
		y := indexRow*KeyHeight + 1
		x := 1
		for indexKey, key := range row.Keys {
			layout.Rows[indexRow].Keys[indexKey].X = x
			layout.Rows[indexRow].Keys[indexKey].Y = y

			x += key.GetWidth()
		}
	}
}

func (layout *Layout) Reset() {
	for rowIndex, row := range layout.Rows {
		for keyIndex := range row.Keys {
			layout.Rows[rowIndex].Keys[keyIndex].Hold = false
		}
	}
}
