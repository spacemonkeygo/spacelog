// Copyright (C) 2013 Space Monkey, Inc.

package log

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func CaptureOutputToFd(fd int) error {
	err := syscall.Dup2(fd, syscall.Stdout)
	if err != nil {
		return err
	}
	err = syscall.Dup2(fd, syscall.Stderr)
	if err != nil {
		return err
	}
	return nil
}

func CaptureOutputToFile(path string) error {
	fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()
	return CaptureOutputToFd(int(fh.Fd()))
}

func CaptureOutputToProcess(tag, command string) error {
	cmd := exec.Command(command, "-t", tag)
	out, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer out.Close()
	type fder interface {
		Fd() uintptr
	}
	out_fder, ok := out.(fder)
	if !ok {
		return fmt.Errorf("unable to get underlying pipe")
	}
	err = CaptureOutputToFd(int(out_fder.Fd()))
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("captured output process died! %s", err)
		}
	}()
	return nil
}
