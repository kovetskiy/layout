package main

import "github.com/gdamore/tcell"

func watchAndDraw(layout *Layout, style *Style) {
	var (
		up   = make(chan string)
		down = make(chan string)
	)
	go watchKeyPress(up, down)

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	err = screen.Init()
	if err != nil {
		panic(err)
	}

	addKeyXY(layout)
	drawLayout(layout, screen)

	quit := make(chan struct{})
	go func() {
		for {
			event := screen.PollEvent()
			switch event := event.(type) {
			case *tcell.EventKey:
				if event.Key() == tcell.KeyCtrlC {
					close(quit)
				}
			}
		}
	}()

	var styleTapHold tcell.Style
	styleTapHold = styleTapHold.Background(tcell.Color(style.Tap.Hold.Background))
	styleTapHold = styleTapHold.Foreground(tcell.Color(style.Tap.Hold.Foreground))

	var styleModHold tcell.Style
	styleModHold = styleModHold.Background(tcell.Color(style.Mod.Hold.Background))
	styleModHold = styleModHold.Foreground(tcell.Color(style.Mod.Hold.Foreground))

	keys := getKeysTable(layout)

	screen.Show()
loop:
	for {
		var name string
		select {
		case name = <-up:
			keys[name].Hold = false
		case name = <-down:
			keys[name].Hold = true
		case <-quit:
			break loop
		}

		drawKey(keys[name], screen, styleTapHold, styleModHold)
		screen.Show()
	}

	screen.Fini()
}

func drawLayout(layout *Layout, screen tcell.Screen) {
	for _, row := range layout.Rows {
		for _, key := range row.Keys {
			drawKey(&key, screen, 0, 0)
		}
	}
}

func drawKey(
	key *Key,
	screen tcell.Screen,
	styleTapHold tcell.Style,
	styleModHold tcell.Style,
) {
	lines := key.Render(nil)

	var style tcell.Style
	if key.Hold {
		if key.Mod {
			style = styleModHold
		} else {
			style = styleTapHold
		}
	}

	for y, line := range lines {
		for x, symbol := range []rune(line) {
			screen.SetContent(
				key.X+x,
				key.Y+y,
				rune(symbol),
				nil,
				style,
			)
		}
	}
}

func getKeysTable(layout *Layout) map[string]*Key {
	keys := map[string]*Key{}
	for indexRow, row := range layout.Rows {
		for indexKey, key := range row.Keys {
			keys[key.Name] = &layout.Rows[indexRow].Keys[indexKey]
		}
	}
	return keys
}
