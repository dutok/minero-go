package config

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Config wraps an unexported Map with methods for parsing and type conversion.
type Config struct {
	parsed      bool
	file, input string
	root        Map // Stores key/value pairs
}

// New creates and initializes a new Config.
func New() *Config { return &Config{root: make(Map)} }

// NewFrom creates and initializes a new Config using m as its initial contents.
func NewFrom(m Map) *Config { return &Config{root: m} }

// String returns the pretty-printed contents of this Config's underlying Map
// sorted by key name.
func (c Config) String() string { return c.root.String() }

// Len returns the number of entries in this Config
func (c Config) Len() int { return len(c.root) }

// Get retrieves a pair value by its key name.
func (c Config) Get(k string) string { return c.root[k] }

// Set sets a key/value pair. Setting an existing key overwrites it.
func (c Config) Set(k, v string) { c.root[k] = v }

// Copy copies this Config's underlying Map and returns that copy.
func (c Config) Copy() (m Map) {
	m = make(Map)
	for k, v := range c.root {
		m[k] = v
	}
	return
}

// Parse parses a string into a config.
func (c Config) Parse(s string) error {
	c.file = "<string>"
	c.input = s
	return c.parse()
}

// Parse parses a file's contents into a config.
func (c Config) ParseFile(f string) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return fmt.Errorf("Couldn't read config file %q.", f)
	}
	c.file = f
	c.input = string(buf)
	return c.parse()
}

// Save saves a config into a file.
func (c Config) Save(file string) error {
	var s = SortedMap(c.root)
	return ioutil.WriteFile(file, []byte(s), 0666)
}

func (c Config) parse() error {
	if c.parsed {
		return nil
	}
	if c.input == "" {
		return fmt.Errorf("Empty config file %q.", c.file)
	}

	// Key chain
	var chain []string

	for index, line := range strings.Split(c.input, "\n") {
		var (
			values = strings.SplitN(line, ":", 2)
			key    = strings.TrimSpace(values[0])
		)

		// Empty line
		if len(values) == 1 || strings.HasPrefix(key, "#") {
			continue
		}

		var (
			value = strings.TrimSpace(values[1])
			level = lineLevel(line)
		)

		// Current value defines new map
		if value == "" {
			if len(chain) != level+1 {
				chain = append(chain, key)
			}
			chain[level] = key
			continue
		}

		// Determine where to save that section
		switch level {
		case 0:
			// First line before this switch
		default:
			if !c.hasRoot(chain, level) {
				return fmt.Errorf("%s:%d %q at level %d has no root.", c.file, index, key, level)
			}

			// Join all key strings
			key = strings.Join(append(chain[:level], key), ".")
		}

		c.root[key] = value
	}

	c.parsed = true
	// Let GC know we are done
	c.input = ""

	return nil
}

// hasRoot ensures a chain of sections exists, linking levels [0, end].
func (c Config) hasRoot(chain []string, end int) bool {
	for i := 0; i < end; i++ {
		if chain[i] == "" {
			return false
		}
	}
	return true
}

func lineLevel(s string) (n int) {
	for _, r := range s {
		if r == '\t' || r == ' ' {
			n++
		} else {
			return
		}
	}
	return
}
