package main

import (
	"ghost-msim/msim"
	"ghost-msim/utils"
)

func main() {
	utils.Log("SYS", "Initializing ghost MSIM")
	msim.Initialise(1863)
}
