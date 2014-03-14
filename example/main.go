package main

import (
	"log"

	"github.com/SpaceMonkeyInc/flagfile"

	space_log "code.spacemonkey.com/go/space/log"
)

var (
	logger = space_log.GetLogger()
)

func main() {
	flagfile.Load()
	err := space_log.Setup("test")
	if err != nil {
		panic(err)
	}
	logger.Debug("hello")
	log.Printf("whoaaa")
	logger.Warn("uh")
	logger.Error("uh oh")
}
