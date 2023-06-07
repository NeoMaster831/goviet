package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	osu "pr0j3ct5/goviet/src/osu/parser"
)

var folders []osu.Osu

func loopFolders(path string) {
	files, _ := ioutil.ReadDir(path)

	for _, file := range files {

		newpath := filepath.Join(path, file.Name())
		if file.IsDir() {
			loopFolders(newpath)
		}

		var sample osu.Osu
		if filepath.Ext(newpath) != ".osu" {
			continue
		}

		sample.Parse(newpath)
		folders = append(folders, sample)

	}
}

func main() {
	loopFolders("C:/Users/last_/AppData/Local/osu!/Songs")
	for _, osu := range folders {
		fmt.Println(osu.Artist + " - " + osu.Title + " [" + osu.Version + "] made by " + osu.Creator)
	}
}
