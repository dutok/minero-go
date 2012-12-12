// nbt provides the ability to read/write NBT data structures from Readers and
// Writers.
//
// NBT is the data serialization format used in many places in the official
// Notchian Minecraft server, typically to represent structured world, chunk
// and player information.
//
// An NBT data structure can be created with code such as the following:
//
//   root := &Compound{
//     map[string]Tag{
//       "Data": &Compound{
//         map[string]Tag{
//           "Byte":   &Byte{1},
//           "Short":  &Short{2},
//           "Int":    &Int{3},
//           "Long":   &Long{4},
//           "Float":  &Float{5},
//           "Double": &Double{6},
//           "String": &String{"foo"},
//           "List":   &List{TagByte, []Tag{&Byte{1}, &Byte{2}}},
//         },
//       },
//     },
//   }
//
// It is required that the root structure be a Compound for compatibility with
// existing NBT structures observed in the official server.
//
// NBT structures can be read from an io.Reader with the Read function.
//
// Many thanks to #mcdevs from Freenode and it's great documentation:
// http://wiki.vg/NBT
package nbt
