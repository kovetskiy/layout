package main

import (
	"log"
	"sort"
	"strings"
)

func groupStates(layout *Layout, states []string) []string {
	keysTable := getKeysTable(layout)

	table := map[string]map[string]struct{}{}
	for _, state := range states {
		keys := strings.Split(state, "+")

		mods, taps := splitKeys(keysTable, keys)
		sort.Strings(mods)

		stateMods := strings.Join(mods, "+")
		_, ok := table[stateMods]
		if !ok {
			table[stateMods] = map[string]struct{}{}
		}

		for _, tap := range taps {
			table[stateMods][tap] = struct{}{}
		}
	}

	groups := []string{}
	for mods, stateTaps := range table {
		taps := []string{}
		for tap, _ := range stateTaps {
			taps = append(taps, tap)
		}

		sort.Strings(taps)

		chunks := []string{}
		chunks = append(chunks, mods)
		chunks = append(chunks, taps...)

		groups = append(groups, strings.Join(chunks, "+"))
	}

	return groups
}

func splitKeys(
	table map[string]*Key,
	keys []string,
) (mods []string, taps []string) {
	for _, name := range keys {
		_, ok := table[name]
		if !ok {
			log.Fatalf("unexpected key given: %s", name)
		}

		if table[name].Mod {
			mods = append(mods, name)
		} else {
			taps = append(taps, name)
		}
	}

	return mods, taps
}
