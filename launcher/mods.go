package launcher

import (
	"fmt"
	"os"
	"strings"
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

func CollectMods(path string) ([]Mod, error) {
	mods := []Mod{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			modName := strings.Split(e.Name(), ".")
			mod := Mod{
				Name: strings.Join(modName, "."),
				Path: path,
				Type: getModType(e.Name()),
			}
			mods = append(mods, mod)
		}
	}
	return mods, nil
}

func ApplyMods(mods *[]Mod) error {
	for _, m := range *mods {
		fmt.Printf("Applying %s\n", m.Name)
		switch m.Type {
		case MOD_DLLFILE:
			// inject
		case MOD_PAKFILE:
			// Add pak and find smallest sig
		case MOD_PATCHFILE:
			// patch paks and make backups
		}
	}

	return nil
}
