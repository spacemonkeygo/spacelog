package main

import (
	"flag"
	"fmt"
	"io"
	"runtime"

	"github.com/SpaceMonkeyInc/flagfile"

	space_log "code.spacemonkey.com/go/space/log"
)

var (
	gosched    = flag.Bool("gosched", false, "if true, call gosched before logging")
	gomaxprocs = flag.Int("gomaxprocs", 1, "gomaxprocs")

	logger = space_log.GetLogger()
)

func main() {
	flagfile.Load()
	runtime.GOMAXPROCS(*gomaxprocs)
	space_log.MustSetup("syslog_test")

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; ; i++ {
				if *gosched {
					runtime.Gosched()
				}
				logger.Warn("hello")
			}
		}()
	}

	var val string
	fmt.Printf("started\n")
	for {
		_, err := fmt.Scanln(&val)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("error: %s\n", err)
			continue
		}
		logger.Warn(val)
		fmt.Printf("val: %s\n", val)
	}
}
