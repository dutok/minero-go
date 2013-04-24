package config

import (
	"fmt"
	"sort"
	"strings"
)

// Map holds all key/value pairs within a Config.
type Map map[string]string

func (m Map) String() string {
	return "Map{\n" + PrettyMap(m) + "\n}"
}

// PrettyMap sorts a Map by its keys and returns a pretty-printed string version.
func PrettyMap(m Map) string {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var values []string
	for _, k := range keys {
		values = append(values, fmt.Sprintf("%q: %s", k, m[k]))
	}
	r := strings.Join(values, "\n")
	return fmt.Sprintf("%s", r)
}

func fileOutput(m Map) string {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Not-so-pretty way to keep track of what section were written
	var headers = make(map[string]bool)

	var values []string
	var key string
	for _, k := range keys {
		parts := strings.Split(k, ".")
		dots := strings.Count(k, ".")

		// Write header, if required
		var h string
		switch dots {
		case 0: // Split returns a slice of len > 1
		case 1:
			key = parts[1]
			h = parts[0]
		default:
			key = parts[len(parts)-1]
			h = parts[len(parts)-2]
		}

		var tabs string
		// First time header is seen
		if !headers[h] {
			headers[h] = true
			tabs = strings.Repeat("\t", dots-1)
			values = append(values, fmt.Sprintf("%s%s:", tabs, h))
		}

		// Write kv
		tabs = strings.Repeat("\t", dots)
		values = append(values, fmt.Sprintf("%s%s: %s", tabs, key, m[k]))
	}

	return fmt.Sprintf("%s", strings.Join(values, "\n"))
}
