package player

import (
	mct "github.com/minero/minero/types/minecraft"
)

// JustLoginMetadata return default metadata for players who just got in.
func JustLoginMetadata(name string) mct.Metadata {
	meta := mct.NewMetadata()
	meta.Entries[0] = &mct.EntryByte{0}                  // Actions
	meta.Entries[1] = &mct.EntryShort{300}               // Drowning counter
	meta.Entries[5] = &mct.EntryString{mct.String(name)} // Plate name
	meta.Entries[6] = &mct.EntryByte{1}                  // Show plate
	meta.Entries[8] = &mct.EntryInt{0}                   // Potion effects
	return meta
}
