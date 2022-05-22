package main

import (
	"log"
	"runtime"

	"github.com/brokeyourbike/mountefi/cmd"
)

func main() {
	if runtime.GOOS != "darwin" {
		log.Fatal("We only support unix systems at the moment")
		return
	}

	cmd.Execute()
}
