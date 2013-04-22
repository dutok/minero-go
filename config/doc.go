// Simple configuration file format.
//
// Properties:
//   - Less verbose than JSON.
//   - Simpler than YAML.
//   - Easy parsing.
//   - Indentation based.
//
// Notes on Indentation:
//   - Spaces and tabs are equivalent here. Examples:
//     "  \t" == "\t\t " // true
//     "   " == "\t\t\t" // true
//   - You can mix both, although it's not recomended.
//   - Indentation level is computed using: level = num_tabs + num_spaces.
//
// Example input:
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
// After parsing produces:
//
//   var config = Map{
//      "a.b.c": "2",
//      "a.b.d": "3",
//      "a.e.f": "5",
//      "g.h": "7"
//   }
//
// Online config tester: http://play.golang.org/p/FP9hHDBjnN
package config
