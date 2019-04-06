package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	reState = regexp.MustCompile(`[\+/]`)
)

func readStates(filename string) ([]string, error) {
	if filename == "-" {
		filename = "/dev/stdin"
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	states := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		states = append(states, scanner.Text())
	}

	return states, scanner.Err()
}

func groupStates(layout *Layout, states []string) []string {
	keysTable := getKeysTable(layout)

	table := map[string]map[string]struct{}{}
	for _, state := range states {
		keys := reState.Split(state, -1)

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

		if len(taps) == 0 {
			continue
		}

		sort.Strings(taps)

		chunks := []string{}
		chunks = append(chunks, mods)
		chunks = append(chunks, strings.Join(taps, "/"))

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
			log.Printf("unexpected key given: %s", name)
			continue
		}

		if table[name].Mod {
			mods = append(mods, name)
		} else {
			taps = append(taps, name)
		}
	}

	return mods, taps
}
