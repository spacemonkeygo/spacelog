// Copyright (C) 2013 Space Monkey, Inc.

package spacelog

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// CaptureOutputToFd redirects the current process' stdout and stderr file
// descriptors to the given file descriptor, using the dup2 syscall.
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

// CaptureOutputToFile opens a filehandle using the given path, then calls
// CaptureOutputToFd on the associated filehandle.
func CaptureOutputToFile(path string) error {
	fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()
	return CaptureOutputToFd(int(fh.Fd()))
}

// CaptureOutputToProcess starts a process and using CaptureOutputToFd,
// redirects stdout and stderr to the subprocess' stdin.
// CaptureOutputToProcess expects the subcommand to last the lifetime of the
// process, and if the subprocess dies, will panic.
func CaptureOutputToProcess(command string, args ...string) error {
	cmd := exec.Command(command, args...)
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
			panic(fmt.Errorf("captured output process died! %s", err))
		}
	}()
	return nil
}
