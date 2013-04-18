// Simple configuration file format.
//
// Input:
//
//   a:
//    b:
//     c: 2
//     d: 3
//    e:
//     f: 5
//   g:
//    h: 7
//
// Go:
//
//   var config = []*Section{
//      &Section{"a", "0", 0, []*Section{
//          &Section{"b", "1", 1, []*Section{
//              &Section{"c", "2", 2, nil},
//              &Section{"d", "3", 2, nil},
//          }},
//          &Section{"e", "4", 1, []*Section{
//              &Section{"f", "5", 2, nil},
//          }},
//      }},
//      &Section{"g", "6", 0, []*Section{
//          &Section{"h", "7", 1, nil},
//      }},
//   }
package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Map map[string]string

type Config struct {
	parsed      bool
	file, input string
	root        Map // Stores key/value pairs
}

func New() *Config                   { return &Config{root: make(Map)} }
func (c Config) Len() int            { return len(c.root) }
func (c Config) String() string      { return fmt.Sprintf("%s", c.root) }
func (c Config) Get(k string) string { return c.root[k] }
func (c Config) Set(k, v string)     { c.root[k] = v }

func (c *Config) Parse(s string) error {
	c.file = "<string>"
	c.input = s
	return c.parse()
}

func (c *Config) ParseFile(f string) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("Couldn't read config file %q.", f)
	}
	c.file = f
	c.input = string(buf)
	return c.parse()
}

func (c *Config) parse() error {
	if c.parsed {
		return nil
	}
	if c.input == "" {
		return fmt.Errorf("Empty config file %q.", c.file)
	}

	// Key chain
	var chain []string

	for index, line := range strings.Split(c.input, "\n") {
		values := strings.SplitN(line, ":", 2)
		// Empty line
		if len(values) == 1 {
			continue
		}

		var (
			key   = strings.TrimSpace(values[0])
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
func (c *Config) hasRoot(chain []string, end int) bool {
	for i := 0; i < end; i++ {
		if chain[i] == "" {
			return false
		}
	}
	return true
}

func (c *Config) SaveTo(f string) error {
	s := fmt.Sprintf("%s", c.root)
	return ioutil.WriteFile(f, []byte(s), 0666)
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
