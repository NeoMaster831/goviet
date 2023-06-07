package main

import (
	"fmt"
	osu "pr0j3ct5/goviet/src/osu/parser"
)

func main() {
	var sample osu.Osu
	sample.Parse("C:/Users/last_/AppData/Local/osu!/Songs/366079 saradisk - 168 - 401/saradisk - 168 - 401 (Nozhomi) [Torpedo].osu")
	fmt.Println(sample.ObjsByTime...)
}
