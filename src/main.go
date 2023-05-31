package main

import (
	"fmt"
	"pr0j3ct5/goviet/src/utils"
)

func main() {
	pid, err := utils.GetPID("Notepad.exe")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(pid)
	}
}
