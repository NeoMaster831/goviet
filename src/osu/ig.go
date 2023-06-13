/*
 * Ingame values.
 */

package osu

import (
	"fmt"
	"os"
	"path/filepath"
	"pr0j3ct5/goviet/src/osu/parser"
	"pr0j3ct5/goviet/src/utils"
	"sync"
)

/*
################################### VARIABLE SECTION #########################################
*/

// consts (aob, pointer chain, ...)
var (
	OSU_SONGS_PATH         = "C:/Users/last_/AppData/Local/osu!/Songs"
	GENERAL_WORKERS        = 100
	MAX_SIZE               = 0x24000000
	OFFSET                 = 1
	CIRCLE_CLICK_DURATION  = 17
	SLIDER_CLICK_DURATION  = 20
	SPINNER_CLICK_DURATION = 17

	CUR_BEATMAP_A, CUR_BEATMAP_M = [...]uint8{0x8B, 0x0D, 0x00, 0x00, 0x00, 0x00, 0xBA, 0x00, 0x00, 0x00, 0x00, 0xE8, 0x00, 0x00, 0x00, 0x00, 0x83, 0xF8}, "xx????x????x????xx"
	TIMESTAMP_A, TIMESTAMP_M     = [...]uint8{0x2B, 0x05, 0x00, 0x00, 0x00, 0x00, 0xA3, 0x00, 0x00, 0x00, 0x00, 0xEB, 0x0A}, "xx????x????xx"
	KEYSET_A, KEYSET_M           = [...]uint8{0x89, 0x55, 0xD4, 0x8B, 0x15, 0x00, 0x00, 0x00, 0x00, 0x38, 0x02}, "xxxxx????xx"
)

// storing value (variables)
var (
	Keyset     [2]uintptr // -> two uint8's
	State      uintptr    // -> int
	Mods       uintptr    // -> int
	Timestamp  uintptr    // -> int
	CurBeatmap uintptr    // -> *osu.BeatmapInstance
	CursorPos  uintptr    // int[2]
	Resolution uintptr    // int[2]
)

// storing data
var (
	Beatmaps map[int]parser.Osu = make(map[int]parser.Osu)
)

// create workers
var (
	wg    sync.WaitGroup
	loopj chan string = make(chan string)
	mutex             = &sync.Mutex{}
)

/*
################################### INIT FUNCTIONS SECTION #########################################
*/

func loopFoldersWorker() error {
	for path := range loopj {
		files, _ := os.ReadDir(path)

		for _, file := range files {
			if !file.IsDir() {

				newpath := filepath.Join(path, file.Name())
				var sample parser.Osu
				if filepath.Ext(newpath) != ".osu" {
					continue
				}

				sample.Parse(newpath)
				mutex.Lock()
				Beatmaps[sample.Id] = sample
				mutex.Unlock()
			}
		}
		wg.Done()
	}
	return nil
}

func loopFolders(path string) {
	files, _ := os.ReadDir(path)

	wg.Add(1)
	go func() {
		loopj <- path
	}()

	for _, file := range files {
		if file.IsDir() {
			loopFolders(filepath.Join(path, file.Name()))
		}
	}
}

func store(hSnap uintptr, pat []uint8, mask string, offset uintptr) (uintptr, bool) {
	found := false
	var (
		ret  uintptr
		aobj chan utils.AobTask = make(chan utils.AobTask)
	)
	for w := 1; w <= GENERAL_WORKERS; w++ {
		go utils.AsyncScanWorker(hSnap, &wg, &found, &aobj, &ret)
	}
	utils.AsyncScanRegion(hSnap, pat, mask, 0x0, uint(MAX_SIZE), &aobj, &wg)

	if !found {
		return 0, false
	}

	_ = utils.RPM(hSnap, ret+offset, &ret)
	ret &= 0x00000000FFFFFFFF // since we use 32 bit game, we should cur off the first 8 bytes
	return ret, true
}

func InitData(hSnap uintptr) error {

	fmt.Println("Parsing .osu files...")
	for w := 1; w <= GENERAL_WORKERS; w++ {
		go loopFoldersWorker()
	}
	loopFolders(OSU_SONGS_PATH)
	wg.Wait()

	fmt.Println("Loaded", len(Beatmaps), "Beatmaps")
	fmt.Println("Getting necessary values...")

	CurBeatmap, _ = store(hSnap, CUR_BEATMAP_A[:], CUR_BEATMAP_M, 0x2)
	fmt.Println("Got Current Beatmap... (1/5)")
	Timestamp, _ = store(hSnap, TIMESTAMP_A[:], TIMESTAMP_M, 0x7)
	fmt.Println("Got Timestamp... (2/5)")
	State = Timestamp + 0x1F4
	fmt.Println("Got State... (3/5)")
	Mods = Timestamp - 0x414
	fmt.Println("Got Mods... (4/5)")
	//tmp, _ := store(hSnap, KEYSET_A[:], KEYSET_M, 0x5)
	//Keyset[0] = utils.Get32BitPointerChainValue(hSnap, tmp, 0, 8, 0x14)
	//Keyset[1] = Keyset[0] + 0x10
	fmt.Println("Got Keysets.... (5/5)")

	fmt.Printf("CurBeatmap: %x, State: %x, Timestamp: %x, Mods: %x\n", CurBeatmap, State, Timestamp, Mods)
	fmt.Printf("Keyset: %x / %x\n", Keyset[0], Keyset[1])

	return nil
}

/*
################################### DATA GETTING SECTION #########################################
*/

func GetBtmpDueId(id int32) parser.Osu {
	return Beatmaps[int(id)]
}
