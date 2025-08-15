package utils

import (
	"syscall"
	"unsafe"
)

type dword uint32
type long int32
type cBOOL int

type THREADENTRY32 struct {
	dwSize             dword
	cntUsage           dword
	th32ThreadID       dword
	th32OwnerProcessID dword
	tpBasePri          long
	tpDeltaPri         long
	dwFlags            dword
}

var (
	PROCESS_CREATE_THREAD            = uint32(0x0002)
	PROCESS_QUERY_INFORMATION        = uint32(0x0400)
	PROCESS_VM_OPERATION             = uint32(0x0008)
	PROCESS_VM_WRITE                 = uint32(0x0020)
	PROCESS_VM_READ                  = uint32(0x0010)
	PROCESS_ALL_ACCESS               = uint32(0x000F0000) | uint32(0x00100000) | uint32(0xFFFF)
	MEM_COMMIT                       = uint32(0x1000)
	MEM_RESERVE                      = uint32(0x2000)
	PAGE_READWRITE                   = uint32(0x04)
	THREAD_SUSPEND_RESUME            = uintptr(0x0002)
	THREAD_TERMINATE                 = uintptr(0x0001)
	THREAD_QUERY_INFORMATION         = uintptr(0x0040)
	THREAD_QUERY_LIMITED_INFORMATION = uintptr(0x0800)
	THREAD_SET_CONTEXT               = uintptr(0x0010)
	THREAD_SET_INFORMATION           = uintptr(0x0020)
	THREAD_SET_LIMITED_INFORMATION   = uintptr(0x0400)
	THREAD_SET_THREAD_TOKEN          = uintptr(0x0080)
	THREAD_GET_CONTEXT               = uintptr(0x0008)
	THREAD_IMPERSONATE               = uintptr(0x0100)
	THREAD_DIRECT_IMPERSONATION      = uintptr(0x0200)
	THREAD_ALL_ACCESS                = uintptr(0x1F03FF)
	TH32CS_SNAPTHREAD                = uint32(0x00000004)
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
	procThread32First      = kernel32.NewProc("Thread32First")
	procThread32Next       = kernel32.NewProc("Thread32Next")

	ntdll = syscall.NewLazyDLL("ntdll.dll")

	procNtSuspendProcess = ntdll.NewProc("NtSuspendProcess")
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
	if int(ret) == -1 {
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

func Thread32First(hSnapshot syscall.Handle, lpThreadEntry uintptr) (cBOOL, error) {
	ret, _, err := procThread32First.Call(uintptr(hSnapshot), lpThreadEntry)
	if ret == 0 {
		return cBOOL(ret), err
	}
	return cBOOL(ret), nil
}

func Thread32Next(hSnapshot syscall.Handle, lpThreadEntry uintptr) (cBOOL, error) {
	ret, _, err := procThread32Next.Call(uintptr(hSnapshot), lpThreadEntry)
	if ret == 0 {
		return cBOOL(ret), err
	}
	return cBOOL(ret), nil
}

// not really Ntsuspendprocess but its better formatted for Go
func NtSuspendProcess(pID dword) error {
	processHandle, err := syscall.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pID))
	if err != nil {
		return err
	}
	procNtSuspendProcess.Call(uintptr(processHandle))
	if err := syscall.CloseHandle(processHandle); err != nil {
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
