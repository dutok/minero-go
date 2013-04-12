// Package status is used with packet Entity Status (0x26).
package status

const (
	_                    = iota // 0
	_                           // 1
	EntityHurt                  // 2
	EntityDead                  // 3
	_                           // 4
	_                           // 5
	WolfTaming                  // 6
	WolfTamed                   // 7
	WolfShakingWater            // 8
	ServerAcceptedEating        // 9
	SheepEating                 // 10
)
