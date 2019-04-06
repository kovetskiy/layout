package main

import (
	"fmt"
	"log"
	"strings"
)

func printLayoutStates(
	layout *Layout,
	style *Style,
	states []string,
) {
	keys := getKeysTable(layout)

	for i, state := range states {
		if i > 0 {
			fmt.Println()
		}

		fmt.Println(state)

		hold := strings.Split(state, "+")

		for _, name := range hold {
			_, ok := keys[name]
			if !ok {
				log.Fatalf("unexpected key given: %s", name)
			}

			keys[name].Hold = true
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
