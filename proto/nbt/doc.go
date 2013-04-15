// Package nbt implements NBT read/write file support.
//
// Its based heavily onthe io.ReaderFrom and io.WriteTo interfaces, so you can
// just use io.Copy to read and write NBT streams.
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
// Many thanks to #mcdevs from Freenode and it's great documentation:
// http://wiki.vg/NBT
package nbt
