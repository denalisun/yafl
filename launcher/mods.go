package launcher

import (
	"bytes"
	"fmt"
	"io"
	"math"
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

func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func ApplyPaks(mods *[]Mod, inst *utils.YAFLInstance) error {
	for _, m := range *mods {
		if m.Type == MOD_PAKFILE {
			contentPath := filepath.Join(inst.BuildPath, "FortniteGame\\Content\\Paks")
			smallestSig, err := findSmallestSig(inst.BuildPath)
			if err != nil {
				return err
			}
			err = CopyFileContents(m.Path, filepath.Join(contentPath, m.Name))
			if err != nil {
				return err
			}
			err = CopyFileContents(smallestSig, filepath.Join(contentPath, strings.Split(filepath.Base(m.Path), ".")[0]+".sig"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ApplyDLLs(mods *[]Mod, inst *utils.YAFLInstance, fortnitePID utils.DWORD) error {
	for _, m := range *mods {
		if m.Type == MOD_DLLFILE {
			err := injectDll(m.Path, fortnitePID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ApplyPatches(mods *[]Mod, inst *utils.YAFLInstance) ([]string, error) {
	moddedFiles := []string{}
	for _, m := range *mods {
		if m.Type == MOD_PATCHFILE {
			f, err := os.ReadFile(m.Path)
			if err != nil {
				return moddedFiles, err
			}

			if !bytes.Equal(f[:2], []byte{0xAE, 'p'}) {
				return moddedFiles, fmt.Errorf("patch \"%s\" doesn't have a valid magic", m.Name)
			}

			pBytes := f[2:]
			separatedBytes := map[int][]byte{}
			i := 0
			for _, byt := range pBytes {
				if byt == 0xFF {
					i++
				} else {
					if separatedBytes[i] == nil {
						separatedBytes[i] = []byte{}
					}
					separatedBytes[i] = append(separatedBytes[i], byt)
				}
			}

			for _, bytes := range separatedBytes {
				instruction := bytes[0]
				if instruction == 0x05 {
					continue
				}
				size := bytes[1]
				parameter := bytes[2:]
				fmt.Println(instruction, size, parameter)
			}
		}
	}
	return moddedFiles, nil
}

func findSmallestSig(path string) (string, error) {
	contentPath := filepath.Join(path, "FortniteGame\\Content\\Paks")
	files, err := os.ReadDir(contentPath)
	if err != nil {
		return "", err
	}

	var sigPath string
	var sigSize int64 = math.MaxInt64
	for _, f := range files {
		splitName := strings.Split(f.Name(), ".")
		if splitName[len(splitName)-1] == "sig" {
			sigPathTemp := filepath.Join(contentPath, f.Name())
			fs, err := os.Stat(sigPathTemp)
			if err != nil {
				return "", err
			}

			size := fs.Size()
			if size < sigSize {
				sigPath = sigPathTemp
				sigSize = size
			}
		}
	}

	return sigPath, nil
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

func RemoveMods(inst utils.YAFLInstance) error {
	contentPath := filepath.Join(inst.BuildPath, "FortniteGame\\Content\\Paks")
	files, err := os.ReadDir(contentPath)
	if err != nil {
		return err
	}
	for _, f := range files {
		splitName := strings.Split(f.Name(), ".")
		if splitName[len(splitName)-1] == ".bak" {
			realPak := strings.Join(splitName[:len(splitName)-1], ".")
			fmt.Println(realPak)
		}
	}

	return nil
}
