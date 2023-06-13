package utils

import "math"

type Vec2 struct {
	X float64
	Y float64
}

var (
	WINDOWS_SIZE = Vec2{X: 2560, Y: 1600}
)

func Dist(from Vec2, to Vec2) float64 {
	return math.Sqrt(math.Pow(from.X-to.X, 2) + math.Pow(from.Y-to.Y, 2))
}

// top left is the origin
func GetPlayfieldSzNOrig() (size Vec2, orig Vec2, playratio float64) {
	ratio := WINDOWS_SIZE.Y / 480
	size = Vec2{X: 512 * ratio, Y: 384 * ratio}
	orig = Vec2{X: (WINDOWS_SIZE.X - size.X) / 2, Y: (WINDOWS_SIZE.Y-size.Y)/4*3 + (-16 * ratio)}
	playratio = size.Y / 384
	return
}
