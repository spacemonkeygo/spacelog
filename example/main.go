package main

import (
	"flag"
	"log"

	"github.com/SpaceMonkeyInc/flagfile"

	space_log "code.spacemonkey.com/go/space/log"
)

var (
	skip_setup = flag.Bool("skip_setup", false, "if true, skip space_log setup")
	logger     = space_log.GetLogger()
)

func main() {
	flagfile.Load()
	if !*skip_setup {
		err := space_log.Setup("test")
		if err != nil {
			panic(err)
		}
	}
	logger.Debug("hello")
	log.Printf("whoaaa")
	logger.Warn("uh")
	logger.Error("uh oh")
}
