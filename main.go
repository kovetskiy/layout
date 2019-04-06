package main

import (
	"log"

	"github.com/docopt/docopt-go"
)

var (
	version = "[manual build]"
	usage   = "layout " + version + `


Usage:
  layout [options] -P [-k]
  layout [options] -P [-m] [<state>...]
  layout [options] -P [-m] -f <file>
  layout [options] -W [-f <file>]
  layout -h | --help
  layout --version

Options:
  -W                       Watch and draw key presses.
  -P                       Print layout and highlight pressed keys.
  -k                       Print only layout keys.
  -m                       Print grouped by mods.
  -f <file>                Filter presses by given file.
  -h --help                Show this screen.
  --style-tap-hold-bg <n>  Background of tap key when hold [default: 11].
  --style-tap-hold-fg <n>  Foreground of tap key when hold [default: 0].
  --style-mod-hold-bg <n>  Background of mod key when hold [default: 1].
  --style-mod-hold-fg <n>  Foreground of mod key when hold [default: 0].
  --version                Show version.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	style, err := getStyle(args)
	if err != nil {
		log.Fatalln(err)
	}

	layout := DefaultLayout

	states, _ := args["<state>"].([]string)
	if filename, ok := args["-f"].(string); ok {
		states, err = readStates(filename)
		if err != nil {
			log.Fatalln(err)
		}
	}

	switch {
	case args["-W"].(bool):
		watchAndDraw(layout, style, states)

	case args["-P"].(bool):
		if args["-k"].(bool) {
			printLayoutKeys(layout)
			return
		}

		fallthrough

	default:
		if args["-m"].(bool) {
			states = groupStates(layout, states)
		}

		printLayoutStates(
			layout,
			style,
			states,
		)
	}
}
