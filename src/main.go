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

	for {

		var (
			beatmap   parser.BeatmapInstance
			timestamp int32
			state     int32
			mods      int32
		)

		bmInstance, _ := utils.Get32BitPtr(hSnap, osu.CurBeatmap)
		utils.RPM(hSnap, bmInstance, &beatmap)
		utils.RPM(hSnap, osu.Timestamp, &timestamp)
		utils.RPM(hSnap, osu.State, &state)
		utils.RPM(hSnap, osu.Mods, &mods)

		fmt.Println(beatmap.SetId, "-", beatmap.Id, timestamp, state, mods)

		time.Sleep(1000 * time.Millisecond)
	}
}
