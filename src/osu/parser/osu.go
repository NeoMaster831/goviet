/*
* .osu Parser. Just implemented important modules
 */

package osu

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Timing Point Parsing is ez enough
type TimingPoint struct {
	Time        int
	BeatLength  float64
	Meter       int
	Uninherited int // 0 or 1
}

const (
	HIT_CIRCLE = 1 << 0
	SLIDER     = 1 << 1
	SPINNER    = 1 << 3
)

// This is pretty hard to optimize, huh
type HitObject struct {
	X, Y         int
	Time         int
	Type         int // 0x00000000
	ObjectParams interface{}
}

type SliderParams struct {
	CurveType   string
	CurvePoints [][2]int
	Slides      int
	Length      float64
}

type SpinnerParams struct {
	EndTime int // End Time, based on BEATMAP'S AUDIO.
}

type Osu struct {

	/* Element makes .osu files distinctly */
	Id int // Beatmap ID

	/* Elements used to calculate something */
	Mode       int           // Mode. 0 = osu!standard, 1 = osu!taiko, 2 = ...
	HP         float64       // HP drain rate
	CS         float64       // Circle Size
	OD         float64       // Overall Difficulty
	AR         float64       // Approach Rate
	SM         float64       // Slider Multiplier
	ST         float64       // Slider tick rate
	ObjsByTime []interface{} // A vector contains Objects(Timing Point, Hit Object) sorted in chronological order

	/* Elements used to show something in user interface */
	Title   string // Song title (Romanised)
	Artist  string // Artist
	Creator string // Beatmap Creator
	Version string // Difficulty Name

}

func parseTimingPoint(form string) TimingPoint {

	var ret TimingPoint
	elems := strings.Split(form, ",")

	ret.Time, _ = strconv.Atoi(elems[0])
	ret.BeatLength, _ = strconv.ParseFloat(elems[1], 64)
	ret.Meter, _ = strconv.Atoi(elems[2])
	ret.Uninherited, _ = strconv.Atoi(elems[6])

	return ret
}

func parseHitObject(form string) HitObject {

	var ret HitObject
	elems := strings.Split(form, ",")

	ret.X, _ = strconv.Atoi(elems[0])
	ret.Y, _ = strconv.Atoi(elems[1])
	ret.Time, _ = strconv.Atoi(elems[2])
	ret.Type, _ = strconv.Atoi(elems[3])

	if ret.Type&SPINNER != 0 { // Case: Spinner
		endTime, _ := strconv.Atoi(elems[5])
		ret.ObjectParams = SpinnerParams{endTime}
	} else if ret.Type&SLIDER != 0 { // Case: Slider
		var apply SliderParams
		vals := strings.Split(elems[5], "|")
		apply.CurveType = vals[0]

		for _, val := range vals[1:] {
			points := strings.Split(val, ":")
			x, _ := strconv.Atoi(points[0])
			y, _ := strconv.Atoi(points[1])
			apply.CurvePoints = append(apply.CurvePoints, [2]int{x, y})
		}

		apply.Slides, _ = strconv.Atoi(elems[6])
		apply.Length, _ = strconv.ParseFloat(elems[7], 64)
		ret.ObjectParams = apply
	}

	return ret
}

func (osu *Osu) Parse(path string) {

	inpf, _ := os.Open(path)
	defer inpf.Close()

	reader := bufio.NewReader(inpf)
	mode := "G" // "G" for general, "T" for timing point section, "H" for hit object section

	var (
		tq []TimingPoint
		hq []HitObject
	)

	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix || err != nil {
			break
		}
		if len(line) == 0 {
			continue
		}

		lstring := string(line[:])

		switch {

		// [General]
		case strings.Contains(lstring, "Mode"):
			osu.Mode, _ = strconv.Atoi(strings.Split(lstring, ": ")[1])

		// [Metadata]
		case strings.Contains(lstring, "Title:"):
			osu.Title = strings.Join(strings.Split(lstring, ":")[1:], "")
		case strings.Contains(lstring, "Artist:"):
			osu.Artist = strings.Join(strings.Split(lstring, ":")[1:], "")
		case strings.Contains(lstring, "Creator:"):
			osu.Creator = strings.Join(strings.Split(lstring, ":")[1:], "")
		case strings.Contains(lstring, "Version:"):
			osu.Version = strings.Join(strings.Split(lstring, ":")[1:], "")
		case strings.Contains(lstring, "BeatmapID:"):
			osu.Id, _ = strconv.Atoi(strings.Split(lstring, ":")[1])

		// [Difficulty]
		case strings.Contains(lstring, "HPDrainRate:"):
			osu.HP, _ = strconv.ParseFloat(strings.Split(lstring, ":")[1], 64)
		case strings.Contains(lstring, "CircleSize:"):
			osu.CS, _ = strconv.ParseFloat(strings.Split(lstring, ":")[1], 64)
		case strings.Contains(lstring, "OverallDifficulty:"):
			osu.OD, _ = strconv.ParseFloat(strings.Split(lstring, ":")[1], 64)
		case strings.Contains(lstring, "ApproachRate:"):
			osu.AR, _ = strconv.ParseFloat(strings.Split(lstring, ":")[1], 64)
		case strings.Contains(lstring, "SliderMultiplier:"):
			osu.SM, _ = strconv.ParseFloat(strings.Split(lstring, ":")[1], 64)
		case strings.Contains(lstring, "SliderTickRate:"):
			osu.ST, _ = strconv.ParseFloat(strings.Split(lstring, ":")[1], 64)

		}

		if lstring == "[TimingPoints]" {
			mode = "T"
			continue
		} else if lstring == "[HitObjects]" {
			mode = "H"
			continue
		} else if strings.Contains(lstring, "[") && strings.Contains(lstring, "]") {
			mode = "G"
			continue
		}

		if mode == "T" {
			tq = append(tq, parseTimingPoint(lstring))
		}
		if mode == "H" {
			hq = append(hq, parseHitObject(lstring))
		}
	}

	// Algorithm: Two pointer distortion
	tq = append(tq, TimingPoint{Time: 2147483647})
	hq = append(hq, HitObject{Time: 2147483647})
	l, r := 0, 0

	// The Hit object is effected first
	for l < len(tq)-1 && r < len(hq)-1 {
		if tq[l].Time < hq[r].Time {
			osu.ObjsByTime = append(osu.ObjsByTime, tq[l])
			l++
		} else {
			osu.ObjsByTime = append(osu.ObjsByTime, hq[r])
			r++
		}
	}

}
