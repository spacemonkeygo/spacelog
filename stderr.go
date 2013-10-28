// Copyright (C) 2013 Space Monkey, Inc.

package log

import (
    "flag"
    "log"
    "os"
    "os/exec"
    "syscall"

    "code.spacemonkey.com/go/errors"
)

var (
    syslog_binary = flag.String("logging_subprocess", "/usr/bin/logger",
        "process to run for stderr-captured logging")
    OutputCaptureError = errors.New(nil, "output capture error")
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

func CaptureOutputToSyslog(tag string) error {
    cmd := exec.Command(*syslog_binary, "-t", tag)
    out, err := cmd.StdinPipe()
    if err != nil {
        return err
    }
    defer out.Close()
    out_fh, ok := out.(*os.File)
    if !ok {
        return OutputCaptureError.New("unable to get underlying pipe")
    }
    err = CaptureOutputToFd(int(out_fh.Fd()))
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
