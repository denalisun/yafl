package utils

import (
	"fmt"
	"os"
	"path"
	"syscall"
	"unsafe"
)

func suspendProcess(pID uint32) error {
	hThreadSnapshot, err := syscall.CreateToolhelp32Snapshot(TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return err
	}

	var threadEntry THREADENTRY32
	threadEntry.dwSize = dword(unsafe.Sizeof(THREADENTRY32{}))

	_, err = Thread32First(hThreadSnapshot, uintptr(unsafe.Pointer(&threadEntry)))
	if err != nil {
		return err
	}

	for {
		if ret, _ := Thread32Next(hThreadSnapshot, uintptr(unsafe.Pointer(&threadEntry))); ret == 1 {
			if threadEntry.th32OwnerProcessID == dword(pID) {
				hThread, err := OpenThread(uint32(THREAD_SUSPEND_RESUME)|uint32(THREAD_QUERY_INFORMATION)|uint32(THREAD_GET_CONTEXT), false, int(threadEntry.th32ThreadID))
				if err != nil {
					return err
				}
				_, err = SuspendThread(hThread)
				if err != nil {
					return err
				}
				err = syscall.CloseHandle(hThread)
				if err != nil {
					return err
				}
			}
		} else {
			break
		}
	}

	syscall.CloseHandle(hThreadSnapshot)

	return nil
}

func LaunchInstance(instance YAFLInstance) error {
	fmt.Printf("Launching %s...\n", instance.Name)

	binariesPath := path.Join(instance.BuildPath, "FortniteGame\\Binaries\\Win64")
	shippingPath := path.Join(binariesPath, "FortniteClient-Win64-Shipping.exe")
	eacShippingPath := path.Join(binariesPath, "FortniteClient-Win64-Shipping_EAC.exe")
	launcherPath := path.Join(binariesPath, "FortniteLauncher.exe")

	launchArgs := []string{}

	var attr = os.ProcAttr{
		Dir: binariesPath,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}
	launcherProcess, err := os.StartProcess(launcherPath, launchArgs, &attr)
	if err != nil {
		return err
	}
	if err := NtSuspendProcess(dword(launcherProcess.Pid)); err != nil {
		return err
	}

	eacProcess, err := os.StartProcess(eacShippingPath, launchArgs, &attr)
	if err != nil {
		return err
	}
	if err := NtSuspendProcess(dword(eacProcess.Pid)); err != nil {
		return err
	}

	_, err = os.StartProcess(shippingPath, launchArgs, &attr)
	if err != nil {
		return err
	}

	return nil
}
