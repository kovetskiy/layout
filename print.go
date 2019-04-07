package main

import (
	"fmt"
	"log"
)

func printLayoutStates(
	layout *Layout,
	layoutKeys map[string]*Key,
	style *Style,
	states []string,
) {
	for i, state := range states {
		if i > 0 {
			fmt.Println()
		}

		fmt.Println(state)

		hold := reState.Split(state, -1)

		for _, name := range hold {
			_, ok := layoutKeys[name]
			if !ok {
				log.Printf("unexpected key given: %s", name)
				continue
			}

			layoutKeys[name].Hold = true
		}

		printLayout(layout, style)

		layout.Reset()
	}

	if len(states) == 0 {
		printLayout(layout, style)
	}
}

func printLayout(layout *Layout, style *Style) {
	for _, row := range layout.Rows {
		rowLines := [KeyHeight]string{}
		for _, key := range row.Keys {
			keyLines := key.Render(style)
			for i := 0; i < KeyHeight; i++ {
				rowLines[i] += keyLines[i]
			}
		}

		for i := 0; i < KeyHeight; i++ {
			fmt.Println(rowLines[i])
		}
	}
}

func printLayoutKeys(layout *Layout) {
	for _, row := range layout.Rows {
		for _, key := range row.Keys {
			fmt.Println(key.Name)
		}
	}
}
