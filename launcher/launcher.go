package launcher

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/denalisun/yafl/utils"
)

func LaunchInstance(instance utils.YAFLInstance) (*os.Process, *os.Process, *os.Process, error) {
	fmt.Printf("Launching %s...\n", instance.Name)

	binariesPath := filepath.Join(instance.BuildPath, "FortniteGame\\Binaries\\Win64")
	shippingPath := filepath.Join(binariesPath, "FortniteClient-Win64-Shipping.exe")
	eacShippingPath := filepath.Join(binariesPath, "FortniteClient-Win64-Shipping_EAC.exe")
	launcherPath := filepath.Join(binariesPath, "FortniteLauncher.exe")

	launchArgs := []string{
		"-epicapp=Fortnite",
		"-epicenv=Prod",
		"-epiclocale=en-us",
		"-epicportal",
		"-skippatchcheck",
		"-NOSSLPINNING",
		"-nobe",
		"-fromfl=eac",
		"-fltoken=3db3ba5dcbd2e16703f3978d", "-caldera=eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiYmU5ZGE1YzJmYmVhNDQwN2IyZjQwZWJhYWQ4NTlhZDQiLCJnZW5lcmF0ZWQiOjE2Mzg3MTcyNzgsImNhbGRlcmFHdWlkIjoiMzgxMGI4NjMtMmE2NS00NDU3LTliNTgtNGRhYjNiNDgyYTg2IiwiYWNQcm92aWRlciI6IkVhc3lBbnRpQ2hlYXQiLCJub3RlcyI6IiIsImZhbGxiYWNrIjpmYWxzZX0.VAWQB67RTxhiWOxx7DBjnzDnXyyEnX7OljJm-j2d88G_WgwQ9wrE6lwMEHZHjBd1ISJdUO1UVUqkfLdU5nofBQ",
	}

	attr := os.ProcAttr{
		Dir: binariesPath,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}

	launcherProcess, err := os.StartProcess(launcherPath, append([]string{launcherPath}, launchArgs...), &attr)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := utils.NtSuspendProcess(utils.DWORD(launcherProcess.Pid)); err != nil {
		return nil, nil, nil, err
	}

	eacProcess, err := os.StartProcess(eacShippingPath, append([]string{eacShippingPath}, launchArgs...), &attr)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := utils.NtSuspendProcess(utils.DWORD(eacProcess.Pid)); err != nil {
		return nil, nil, nil, err
	}

	shippingProcess, err := os.StartProcess(shippingPath, append([]string{shippingPath}, launchArgs...), &attr)
	if err != nil {
		return nil, nil, nil, err
	}

	return shippingProcess, launcherProcess, eacProcess, nil
}
