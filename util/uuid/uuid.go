// Package UUID implements UUIDv4 based on RFC 4122.
//
// Source:
// http://www.ietf.org/rfc/rfc4122.txt
package uuid

import (
	"fmt"
	"math/rand"
)

// type UUID struct {
// 	time_low                  [4]byte // 0-3    The low field of the timestamp
// 	time_mid                  [2]byte // 4-5    The middle field of the timestamp
// 	time_hi_and_version       [2]byte // 6-7    The high field of the timestamp multiplexed with the version number
// 	clock_seq_hi_and_reserved byte    // 8      The high field of the clock sequence multiplexed with the variant
// 	clock_seq_low             byte    // 9      The low field of the clock sequence
// 	node                      [6]byte // 10-15  The spatially unique node identifier
// }

type UUID [16]byte

func (u UUID) String() string {
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		u[0], u[1], u[2], u[3],
		u[4], u[5],
		u[6], u[7],
		u[8], u[9],
		u[10], u[11], u[12], u[13], u[14], u[15])
}

// UUID4 generates a v4 UUID using pseudo-random numbers.
func UUID4() (u UUID) {
	// s := make([]byte, 16)
	// n, err := rand.Read(s)
	// if n != 16 || err != nil {
	// 	return nil
	// }

	for i, _ := range u {
		// u[i] = v
		u[i] = byte(rand.Int31n(256))
	}

	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number from Section 4.1.3.
	u[6] = u[6]&0x0f | 0x40

	// Set the two most significant bits (bits 6 and 7) of the
	// clock_seq_hi_and_reserved to zero and one, respectively.
	u[8] = u[8]&0x3f | 0x80

	return u
}
