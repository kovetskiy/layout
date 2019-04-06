package main

import (
	"strconv"

	"github.com/reconquest/karma-go"
)

type Style struct {
	Tap struct {
		Hold struct {
			Background int
			Foreground int
		}
	}
	Mod struct {
		Hold struct {
			Background int
			Foreground int
		}
	}
}

func getStyle(args map[string]interface{}) (*Style, error) {
	var err error
	get := func(name string) (int, error) {
		arg := "--style-" + name
		value, err := strconv.Atoi(args[arg].(string))
		if err != nil {
			return 0, karma.Format(
				err,
				"unexpected value of %q: %q",
				arg, args[arg].(string),
			)
		}

		return value, nil
	}

	var style Style
	style.Tap.Hold.Background, err = get("tap-hold-bg")
	if err != nil {
		return nil, err
	}

	style.Tap.Hold.Foreground, err = get("tap-hold-fg")
	if err != nil {
		return nil, err
	}

	style.Mod.Hold.Background, err = get("mod-hold-bg")
	if err != nil {
		return nil, err
	}

	style.Mod.Hold.Foreground, err = get("mod-hold-fg")
	if err != nil {
		return nil, err
	}

	return &style, nil
}
