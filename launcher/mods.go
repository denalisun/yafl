package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/denalisun/yafl/utils"
	"golang.org/x/sys/windows"
)

const (
	MOD_PAKFILE int = iota
	MOD_DLLFILE
	MOD_PATCHFILE
	MOD_INVALID int = -1
)

type Mod struct {
	Name string
	Path string
	Type int
}

func (m *Mod) apply(pID utils.DWORD, inst *utils.YAFLInstance) error {
	fmt.Printf("Applying %s...\n", m.Name)
	switch m.Type {
	case MOD_DLLFILE:
		err := injectDll(m.Path, pID)
		if err != nil {
			return err
		}
	case MOD_PAKFILE:
		// Add pak and find smallest sig
	case MOD_PATCHFILE:
		// patch paks and make backups
	}

	fmt.Printf("Successfully applied %s!\n", m.Name)

	return nil
}

func getModType(name string) int {
	split := strings.Split(name, ".")
	switch split[len(split)-1] {
	case "pak":
		return MOD_PAKFILE
	case "dll":
		return MOD_DLLFILE
	case "bin":
		return MOD_PATCHFILE
	}
	return MOD_INVALID
}

func injectDll(path string, fortnitePID utils.DWORD) error {
	hProcess, err := syscall.OpenProcess(utils.PROCESS_CREATE_THREAD|utils.PROCESS_QUERY_INFORMATION|utils.PROCESS_VM_OPERATION|utils.PROCESS_VM_WRITE|utils.PROCESS_VM_READ, false, uint32(fortnitePID))
	if err != nil {
		return err
	}
	if hProcess == 0 {
		return fmt.Errorf("failed to open Fortnite process")
	}
	defer syscall.CloseHandle(hProcess)

	addr, err := utils.VirtualAllocEx(hProcess, 0, uint32(len(path)+1), utils.MEM_COMMIT|utils.MEM_RESERVE, utils.PAGE_READWRITE)
	if err != nil {
		return err
	}
	if addr == 0 {
		return fmt.Errorf("failed to allocate memory in Fortnite process")
	}

	pathAsBytes, err := windows.ByteSliceFromString(path)
	if err != nil {
		return err
	}

	wpmRet, err := utils.WriteProcessMemory(hProcess, addr, pathAsBytes, uint32(len(path)))
	if wpmRet == 0 {
		return err
	}

	llAddr, err := syscall.GetProcAddress(syscall.Handle(utils.Kernel32.Handle()), "LoadLibraryA")
	if err != nil {
		return err
	}

	tHandle, _ := utils.CreateRemoteThread(hProcess, 0, 0, llAddr, addr, 0, 0)
	defer syscall.CloseHandle(tHandle)

	return nil
}

func CollectMods(path string) ([]Mod, error) {
	mods := []Mod{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		modName := strings.Split(e.Name(), ".")
		mod := Mod{
			Name: strings.Join(modName, "."),
			Path: filepath.Join(path, e.Name()),
			Type: getModType(e.Name()),
		}
		mods = append(mods, mod)
	}
	return mods, nil
}

func ApplyMods(mods *[]Mod, fortnitePID utils.DWORD, inst *utils.YAFLInstance) error {
	for _, m := range *mods {
		err := m.apply(fortnitePID, inst)
		if err != nil {
			fmt.Printf("Failed to apply mod: %s\n", err)
			continue
		}
	}

	return nil
}
