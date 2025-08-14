package utils

import (
	"syscall"
	"unsafe"
)

var (
	PROCESS_CREATE_THREAD     = uint32(0x0002)
	PROCESS_QUERY_INFORMATION = uint32(0x0400)
	PROCESS_VM_OPERATION      = uint32(0x0008)
	PROCESS_VM_WRITE          = uint32(0x0020)
	PROCESS_VM_READ           = uint32(0x0010)
	MEM_COMMIT                = uint32(0x1000)
	MEM_RESERVE               = uint32(0x2000)
	PAGE_READWRITE            = uint32(0x04)
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procVirtualAllocEx     = kernel32.NewProc("VirtualAllocEx")
	procWriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
	procCreateRemoteThread = kernel32.NewProc("CreateRemoteThread")
	procGetModuleHandle    = kernel32.NewProc("GetModuleHandle")
	procSuspendThread      = kernel32.NewProc("SuspendThread")
	procResumeThread       = kernel32.NewProc("ResumeThread")
	procOpenThread         = kernel32.NewProc("OpenThread")
	procAllocConsole       = kernel32.NewProc("AllocConsole")
)

// OpenProcess
// GetProcAddress
// CloseHandle

func VirtualAllocEx(hProcess syscall.Handle, lpAddress uintptr, dwSize uint32, flAllocationType uint32, flProtect uint32) (uintptr, error) {
	ret, _, err := procVirtualAllocEx.Call(
		uintptr(hProcess),
		lpAddress,
		uintptr(dwSize),
		uintptr(flAllocationType),
		uintptr(flProtect),
	)
	if ret == 0 {
		return 0, err
	}
	return ret, nil
}

func WriteProcessMemory(hProcess syscall.Handle, lpBaseAddress uintptr, lpBuffer []byte, nSize uint32) (uint32, error) {
	var bytesWritten uint32
	ret, _, err := procWriteProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		uintptr(unsafe.Pointer(&lpBuffer[0])),
		uintptr(nSize),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	if ret == 0 {
		return 0, err
	}
	return bytesWritten, nil
}

func CreateRemoteThread(hProcess syscall.Handle, lpThreadAttributes uintptr, dwStackSize uint32, lpStartAddress uintptr, lpParameter uintptr, dwCreationFlags uint32, lpThreadId uintptr) (syscall.Handle, error) {
	ret, _, err := procCreateRemoteThread.Call(
		uintptr(hProcess),
		lpThreadAttributes,
		uintptr(dwStackSize),
		lpStartAddress,
		lpParameter,
		uintptr(dwCreationFlags),
		lpThreadId,
	)
	if ret == 0 {
		return 0, err
	}
	return syscall.Handle(ret), nil
}

func GetModuleHandle(lpModuleName string) (syscall.Handle, error) {
	ret, _, err := procGetModuleHandle.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpModuleName))),
	)
	if ret == 0 {
		return 0, err
	}
	return syscall.Handle(ret), nil
}

func SuspendThread(hThread syscall.Handle) (uint32, error) {
	ret, _, err := procSuspendThread.Call(uintptr(hThread))
	if ret == 0 {
		return 0, err
	}
	return uint32(ret), nil
}

func ResumeThread(hThread syscall.Handle) (uint32, error) {
	ret, _, err := procResumeThread.Call(uintptr(hThread))
	if ret == 0 {
		return 0, err
	}
	return uint32(ret), nil
}

func OpenThread(dwDesiredAccess uint32, bInheritHandle bool, dwThreadId int) (syscall.Handle, error) {
	ret, _, err := procOpenThread.Call(
		uintptr(dwDesiredAccess),
		uintptr(boolToUint(bInheritHandle)),
		uintptr(dwThreadId),
	)
	if ret == 0 {
		return 0, err
	}
	return syscall.Handle(ret), nil
}

func AllocConsole() error {
	ret, _, err := procAllocConsole.Call()
	if ret == 0 {
		return err
	}
	return nil
}

func boolToUint(b bool) uint {
	if b {
		return 1
	}
	return 0
}
