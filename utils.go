package main

import (
	"syscall"
	"fmt"
)

var (
	ntdll 				 = syscall.NewLazyDLL("ntdll.dll")
	procNtSuspendProcess = ntdll.NewProc("NtSuspendProcess")
)

func NtSuspendProcess(processHandle syscall.Handle) error {
	ret, _, err := procNtSuspendProcess.Call(uintptr(processHandle))
	if ret != 0 {
		return fmt.Errorf("failed to suspend process: %v", err)
	}
	return nil
}
