package main

import (
	"plugins/pods/lib"
	"plugins/util"
)

func main() {
	defer util.ResetSTTY()
	lib.RunCmd()
}
