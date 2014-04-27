// Copyright (C) 2014 Space Monkey, Inc.

package main

import (
	"flag"
	"fmt"
	"io"
	"log/syslog"
	"runtime"

	"github.com/SpaceMonkeyGo/flagfile"

	space_log "code.spacemonkey.com/go/space/log"
)

var (
	gosched    = flag.Bool("gosched", false, "if true, call gosched before logging")
	gomaxprocs = flag.Int("gomaxprocs", 1, "gomaxprocs")
	facility   = flag.Int("facility", int(syslog.LOG_LOCAL0), "syslog facility")

	logger = space_log.GetLogger()
)

func main() {
	flagfile.Load()
	runtime.GOMAXPROCS(*gomaxprocs)
	space_log.MustSetupWithFacility("syslog_test", syslog.Priority(*facility))

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; ; i++ {
				if *gosched {
					runtime.Gosched()
				}
				logger.Notice("hello")
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
		logger.Notice(val)
		fmt.Printf("val: %s\n", val)
	}
}
