package main

import (
	"plugins/deployment/lib"
	"plugins/util"
)

func main() {
	defer util.ResetSTTY()
	lib.RunCmd()
}
