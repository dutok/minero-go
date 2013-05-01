package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var ErrEmpty = errors.New("Empty config file")

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
		return fmt.Errorf("Couldn't read config file %q", f)
	}
	c.file = f
	c.input = string(buf)
	return c.parse()
}

// Save saves a config into a file.
func (c Config) Save(file string) error {
	var s = fileOutput(c.root)
	return ioutil.WriteFile(file, []byte(s), 0666)
}

func (c Config) parse() error {
	if c.parsed {
		return nil
	}
	if c.input == "" {
		return ErrEmpty
	}

	// Key chain
	var chain []string

	for i, line := range strings.Split(c.input, "\n") {
		// Empty/Comment line
		if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		var (
			level  = lineLevel(line)
			values = strings.SplitN(line, ":", 2)
			key    = strings.TrimSpace(values[0])
			value  = strings.TrimSpace(values[1])
		)

		// Remove comments from value
		if index := strings.Index(value, "#"); index != -1 {
			value = strings.TrimSpace(value[:index])
		}

		// New map marker
		switch value {
		case "":
			if len(chain) != level+1 {
				chain = append(chain, key)
			}
			chain[level] = key
			continue
		case "-":
			// Empty string marker
			value = ""
		}

		// Determine where to save that section
		switch level {
		case 0:
			// First line after this switch
		default:
			if !c.hasRoot(chain, level) {
				return fmt.Errorf("%s:%d %q at level %d has no root", c.file, i, key, level)
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
	if len(chain) < end {
		return false
	}
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
