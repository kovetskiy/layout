package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffStates(t *testing.T) {
	test := assert.New(t)

	unhold, hold := diffStates(
		StateGroup{
			Mods: []string{"ALT"},
			Taps: []string{"A"},
		},
		StateGroup{
			Mods: []string{"META"},
			Taps: []string{"B"},
		},
	)

	test.EqualValues([]string{"ALT", "A"}, unhold)
	test.EqualValues([]string{"META", "B"}, hold)
}

func TestFindState(t *testing.T) {
	test := assert.New(t)

	alt := StateGroup{
		Mods: []string{"ALT"},
		Taps: []string{"A"},
	}
	meta := StateGroup{
		Mods: []string{"META"},
		Taps: []string{"B"},
	}

	mods := map[string]struct{}{
		"META": {},
	}

	state := findState([]StateGroup{alt, meta}, mods)

	test.Equal(meta, state)
}
