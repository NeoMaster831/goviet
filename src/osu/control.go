/*
 * utility functions for some shit
 */

package osu

import (
	"fmt"
	"math/rand"
	"pr0j3ct5/goviet/src/osu/parser"
	"pr0j3ct5/goviet/src/utils"
	"pr0j3ct5/goviet/src/winapi"
	"reflect"
	"time"

	"github.com/micmonay/keybd_event"
)

/*
beatLength (Decimal): This property has two meanings:
For uninherited timing points, the duration of a beat, in milliseconds.
For inherited timing points, a negative inverse slider velocity multiplier, as a percentage.
For example, -50 would make all sliders in this timing section twice as fast as SliderMultiplier.
*/

var valset = map[uint8]int{
	// Alphabets
	81: 16, 87: 17, 69: 18, 82: 19, 84: 20,
	89: 21, 85: 22, 73: 23, 79: 24, 80: 25,
	65: 30, 83: 31, 68: 32, 70: 33, 71: 34,
	72: 35, 74: 36, 75: 37, 76: 38, 90: 44,
	88: 45, 67: 46, 86: 47, 66: 48, 78: 49,
	77: 50,

	// Numbers
	49: 2, 50: 3, 51: 4, 52: 5, 53: 6,
	54: 7, 55: 8, 56: 9, 57: 10, 48: 11,

	// Not supporting others rn
}

// variables

const (
	RELAX_WORKERS = 1
)

var (
	nowstate NowPlayState
	idx      int
)

// not converted (raw data)
type NowPlayState struct {
	BeatLength float64
	OrigSM     float64
	SMMul      float64
}

// utility functions

func GetSliderTime(now NowPlayState, slider parser.SliderParams) int {
	fmt.Println(now.OrigSM, now.SMMul, now.BeatLength, slider.Length, slider.Slides, slider.Length/(now.OrigSM*now.SMMul), now.BeatLength)
	return int(slider.Length/(now.OrigSM*100.0)*(now.SMMul/100.0)*now.BeatLength) * slider.Slides
}

// 'bout relax

func AsyncHitCircleClick(key uint8) {
	kb, _ := keybd_event.NewKeyBonding()
	bending := valset[key]
	if bending == 0 {
		return
	}

	kb.SetKeys(bending)
	presleep := rand.Intn(OFFSET * 2)
	time.Sleep(time.Duration(presleep) * time.Millisecond)

	fmt.Println("Debug: Hit Circle Clicked after", presleep, "amount of milliseconds.")
	kb.Press()
	time.Sleep(time.Duration(CIRCLE_CLICK_DURATION) * time.Millisecond)
	kb.Release()
}

func AsyncSliderClick(key uint8, state NowPlayState, slider parser.SliderParams) {
	kb, _ := keybd_event.NewKeyBonding()
	bending := valset[key]
	if bending == 0 {
		return
	}

	duration := GetSliderTime(state, slider)
	presleep := rand.Intn(OFFSET * 2)
	time.Sleep(time.Duration(presleep) * time.Millisecond)

	fmt.Println("Debug: Slider clicked after", presleep, ", Slider duration:", duration)
	kb.SetKeys(bending)
	kb.Press()
	time.Sleep(time.Duration(duration+SLIDER_CLICK_DURATION) * time.Millisecond)
	kb.Release()
}

func AsyncSpinnerClick(key uint8, spinner parser.SpinnerParams, now int) {
	kb, _ := keybd_event.NewKeyBonding()
	bending := valset[key]
	if bending == 0 {
		return
	}

	duration := spinner.EndTime - now
	fmt.Println("Debug: Spinner clicked! Spinner duration:", duration)
	kb.SetKeys(bending)
	kb.Press()
	time.Sleep(time.Duration(duration+SPINNER_CLICK_DURATION) * time.Millisecond)
	kb.Release()
}

func AsyncRelaxController(hSnap uintptr, playing *bool) {

	winapi.ProcTimeBeginPeriod.Call(uintptr(1))
	var (
		beatmap   parser.BeatmapInstance
		timestamp int32
		state     int32
		//mods      int32
		key1 uint8
		//key2      uint8
	)

	idx = 0
	bmInstance, _ := utils.Get32BitPtr(hSnap, CurBeatmap)
	utils.RPM(hSnap, bmInstance, &beatmap)
	cur := GetBtmpDueId(beatmap.Id)

	queue := cur.ObjsByTime
	nowstate = NowPlayState{OrigSM: cur.SM, SMMul: 100.0}

	fmt.Println("Relax Controller started")
	for idx < len(queue) {
		s := time.Now()
		time.Sleep(time.Millisecond)
		utils.RPM(hSnap, Timestamp, &timestamp)
		utils.RPM(hSnap, State, &state)
		//utils.RPM(hSnap, Keyset[0], key1)
		key1 = 68

		if state != 2 { // Play == 2
			*playing = false
			return
		}

		if reflect.TypeOf(queue[idx]).String() == "parser.TimingPoint" { // Case: Timing Point
			obj := queue[idx].(parser.TimingPoint)
			//fmt.Println(obj.Time, int(timestamp)+OFFSET, int(timestamp)-OFFSET, *playing)
			if obj.Time > int(timestamp)+OFFSET {
				continue
			}
			fmt.Println(idx)

			if obj.Uninherited == 1 {
				nowstate.BeatLength = obj.BeatLength
			} else {
				nowstate.SMMul = -obj.BeatLength
			}
			idx++
			fmt.Println("Timing Point applied")

		} else { // Case: Hit Object
			obj := queue[idx].(parser.HitObject)
			//fmt.Println(obj.Time, int(timestamp)+OFFSET, int(timestamp)-OFFSET, *playing)
			if obj.Time > int(timestamp)+OFFSET {
				continue
			}
			fmt.Println(idx)

			if obj.ObjectParams == nil { // Case: Hit Circle
				go AsyncHitCircleClick(key1)
			} else if reflect.TypeOf(obj.ObjectParams).String() == "parser.SliderParams" { // Case: Slider
				go AsyncSliderClick(key1, nowstate, obj.ObjectParams.(parser.SliderParams))
			} else { // Case: Spinner
				go AsyncSpinnerClick(key1, obj.ObjectParams.(parser.SpinnerParams), int(timestamp))
			}
			idx++
		}

		fmt.Println("el:", time.Since(s))
	}

	*playing = false
}
