package main

import "github.com/gdamore/tcell"
import "log"

func watchAndDraw(
	layout *Layout,
	layoutKeys map[string]*Key,
	style *Style,
	states []StateGroup,
) {
	up, down, err := watchKeyPress()
	if err != nil {
		log.Fatalln(err)
	}

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

	screen.Show()

	if len(states) == 0 {
		drawHold(layoutKeys, screen, style, up, down, quit)
	} else {
		drawStates(states, layoutKeys, screen, style, up, down, quit)
	}

	screen.Fini()
}

func drawStates(
	states []StateGroup,
	layoutKeys map[string]*Key,
	screen tcell.Screen,
	style *Style,
	up, down chan string,
	quit chan struct{},
) {
	styleTapHold, styleModHold := getHoldStyles(style)

	mods := map[string]struct{}{}

	var nowState StateGroup
loop:
	for {
		var name string
		var pressed bool
		select {
		case <-quit:
			break loop

		case name = <-up:
			pressed = false

		case name = <-down:
			pressed = true
		}

		key, ok := layoutKeys[name]
		if !ok {
			continue
		}

		if key.Mod {
			if pressed {
				mods[name] = struct{}{}
			} else {
				delete(mods, name)
			}
		}

		nextState := findState(states, mods)

		unhold, hold := diffStates(nowState, nextState)
		for _, name := range unhold {
			layoutKeys[name].Hold = false
			drawKey(layoutKeys[name], screen, styleTapHold, styleModHold, true)
		}
		for _, name := range hold {
			layoutKeys[name].Hold = true
			drawKey(layoutKeys[name], screen, styleTapHold, styleModHold, true)
		}

		nowState = nextState

		screen.Show()
	}
}

func diffStates(
	prev StateGroup,
	next StateGroup,
) (unhold []string, hold []string) {
	prevAll := prev.All()
	nextAll := next.All()

	for _, prevKey := range prevAll {
		found := false
		for _, nextKey := range nextAll {
			if nextKey == prevKey {
				found = true
				break
			}
		}

		if !found {
			unhold = append(unhold, prevKey)
		}
	}

	for _, nextKey := range nextAll {
		found := false
		for _, prevKey := range prevAll {
			if prevKey == nextKey {
				found = true
				break
			}
		}

		if !found {
			hold = append(hold, nextKey)
		}
	}

	return unhold, hold
}

func findState(states []StateGroup, mods map[string]struct{}) StateGroup {
search:
	for _, state := range states {
		if len(state.Mods) == len(mods) {
			for queryMod := range mods {
				found := false
				for _, stateMod := range state.Mods {
					if stateMod == queryMod {
						found = true
					}
				}
				if !found {
					continue search
				}
			}

			return state
		}
	}

	allMods := []string{}
	for mod := range mods {
		allMods = append(allMods, mod)
	}

	return StateGroup{Mods: allMods}
}

func drawHold(
	layoutKeys map[string]*Key,
	screen tcell.Screen,
	style *Style,
	up, down chan string,
	quit chan struct{},
) {
	styleTapHold, styleModHold := getHoldStyles(style)

loop:
	for {
		var name string
		var hold bool
		select {
		case <-quit:
			break loop
		case name = <-up:
			hold = false
		case name = <-down:
			hold = true
		}

		if _, ok := layoutKeys[name]; !ok {
			log.Printf("Unknown key pressed: %s", name)
			continue
		}

		layoutKeys[name].Hold = hold

		drawKey(layoutKeys[name], screen, styleTapHold, styleModHold, true)
		screen.Show()
	}
}

func drawLayout(layout *Layout, screen tcell.Screen) {
	for rowIndex, row := range layout.Rows {
		for keyIndex, _ := range row.Keys {
			drawKey(
				&layout.Rows[rowIndex].Keys[keyIndex],
				screen,
				0,
				0,
				false,
			)
		}
	}
}

func getHoldStyles(style *Style) (tap, mod tcell.Style) {
	tap = tap.Background(tcell.Color(style.Tap.Hold.Background))
	tap = tap.Foreground(tcell.Color(style.Tap.Hold.Foreground))
	mod = mod.Background(tcell.Color(style.Mod.Hold.Background))
	mod = mod.Foreground(tcell.Color(style.Mod.Hold.Foreground))

	return tap, mod
}

func drawKey(
	key *Key,
	screen tcell.Screen,
	styleTapHold tcell.Style,
	styleModHold tcell.Style,
	skipMargin bool,
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
			if skipMargin {
				if x < key.Margin {
					continue
				}
			}
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
