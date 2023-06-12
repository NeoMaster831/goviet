package main

import (
	"fmt"
	"pr0j3ct5/goviet/src/osu"
	"pr0j3ct5/goviet/src/osu/parser"
	"pr0j3ct5/goviet/src/utils"
	"pr0j3ct5/goviet/src/winapi"
	"time"
)

func main() {

	pid, _ := utils.GetPID("osu!.exe")
	mBase, _ := utils.GetModuleBaseAddress(pid, "osu!.exe")
	fmt.Println("mBase:", mBase)
	hSnap, _ := winapi.OpenProcess(0x001FFFFF, 0, pid)
	osu.InitData(hSnap)

	var (
		playing = false
	)

	for {

		var (
			beatmap   parser.BeatmapInstance
			timestamp int32
			state     int32
			mods      int32
			//key1      uint8
			//key2      uint8
		)

		bmInstance, _ := utils.Get32BitPtr(hSnap, osu.CurBeatmap)
		utils.RPM(hSnap, bmInstance, &beatmap)
		//cur := osu.GetBtmpDueId(beatmap.Id)

		utils.RPM(hSnap, osu.Timestamp, &timestamp)
		utils.RPM(hSnap, osu.State, &state)
		utils.RPM(hSnap, osu.Mods, &mods)
		//utils.RPM(hSnap, osu.Keyset[0], &key1)
		//utils.RPM(hSnap, osu.Keyset[1], &Millisecond)

		if state == 2 && !playing {
			playing = true
			for i := 1; i <= osu.RELAX_WORKERS; i++ {
				go osu.AsyncRelaxController(hSnap, &playing)
				time.Sleep(30 * time.Millisecond)
			}
		}
		//fmt.Printf("%s - %s [%s] by %s\n", cur.Artist, cur.Title, cur.Version, cur.Creator)
		//fmt.Printf("Timestamp: %d | State: %d | Mods: 0b%b | Keyset: %c / %c\n", timestamp, state, mods, key1, key2)

		time.Sleep(1000 * time.Millisecond)
	}

}
