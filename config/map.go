package config

import (
	"fmt"
	"sort"
	"strings"
)

// Map holds all key/value pairs within a Config.
type Map map[string]string

func (m Map) String() string {
	return "Map{\n" + SortedMap(m) + "\n}"
}

// SortedMap sorts a Map by its keys and returns it as a string
func SortedMap(m Map) string {
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
