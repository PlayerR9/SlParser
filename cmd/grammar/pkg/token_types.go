package pkg

import (
	"strings"

	gslc "github.com/PlayerR9/mygo-lib/slices"
)

func DetermineTokenTypes(rules []*Rule) []string {
	var symbols []string

	for _, rule := range rules {
		_, _ = gslc.Merge(&symbols, rule.Symbols())
	}

	buckets := make(map[string][]string, 3)
	buckets["E"] = make([]string, 0)
	buckets["T"] = make([]string, 0)
	buckets["N"] = make([]string, 0)

	for _, rhs := range symbols {
		var loc string

		if strings.HasPrefix(rhs, "N") {
			loc = "N"
		} else if strings.HasPrefix(rhs, "E") {
			loc = "E"
		} else {
			loc = "T"
		}

		prev := buckets[loc]
		prev = append(prev, rhs)
		buckets[loc] = prev
	}

	var idx int

	for _, rhs := range buckets["E"] {
		symbols[idx] = rhs
		idx++
	}

	for _, rhs := range buckets["T"] {
		symbols[idx] = rhs
		idx++
	}

	for _, rhs := range buckets["N"] {
		symbols[idx] = rhs
		idx++
	}

	return symbols
}

func FindLastTerminal(symbols []string) string {
	var last_rhs string

	for _, rhs := range symbols {
		if strings.HasPrefix(rhs, "N") {
			return last_rhs
		}

		last_rhs = rhs
	}

	return last_rhs
}
