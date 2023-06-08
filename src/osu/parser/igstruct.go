/*
 * Parses ingame structures by reverse engineering
 */

package parser

import "pr0j3ct5/goviet/src/utils"

const ( // Player's current state (osu.OsuModes)
	PLAY = 2
)

const ( // (osu_common.Mods)
	HIDDEN     = 1 << 3
	HARDROCK   = 1 << 4
	DOUBLETIME = 1 << 6
	NIGHTCORE  = 1 << 9
	SPEEDUP    = DOUBLETIME | NIGHTCORE
)

type BeatmapInstance struct {
	_     [0x94]byte
	Path  *utils.NString // Offset: 94 [4 bytes]
	_     [0x2C]byte
	Id    int32 // Offset: C4 [4 bytes]
	SetId int32 // Offset: C8 [4 bytes]
}
