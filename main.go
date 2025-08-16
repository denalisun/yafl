package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/denalisun/yafl/launcher"
	"github.com/denalisun/yafl/utils"
)

func main() {
	opt, err := utils.ProcessArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := utils.GetData()
	if err != nil {
		fmt.Println(err)
	}

	//TODO: Replace with switch statement

	switch opt.MainOperation {
	case "profile":
		switch opt.SubOperation {
		case "add":
			if len(opt.Parameters) < 2 {
				fmt.Printf("Wrong parameter count! %d provided, 2 or more required!\n", len(opt.Parameters))
				break
			}

			err := utils.CreateInstance(&data, opt.Parameters[0], opt.Parameters[1])
			if err != nil {
				fmt.Println(err)
				break
			}
		case "remove":
			if len(opt.Parameters) != 1 {
				fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
				break
			}
			utils.RemoveInstance(&data, opt.Parameters[0])
		case "list":
			allInstancesFormat := []string{}
			for _, v := range data {
				allInstancesFormat = append(allInstancesFormat, fmt.Sprintf("\t- %s (%s)", v.Name, v.BuildPath))
			}
			fmt.Printf("All instances (%d):\n%s", len(allInstancesFormat), strings.Join(allInstancesFormat, "\n"))
		}
	case "play":
		if len(opt.Parameters) != 1 {
			fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
			break
		}
		inst := utils.FetchInstance(&data, opt.Parameters[0])

		allMods, err := launcher.CollectMods(inst.ModsPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(allMods)

		err = launcher.ApplyMods(&allMods, 0, inst, launcher.MOD_PAKFILE)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = launcher.ApplyMods(&allMods, 0, inst, launcher.MOD_PATCHFILE)
		if err != nil {
			fmt.Println(err)
			return
		}

		shippingProcess, launcherProcess, eacProcess, err := launcher.LaunchInstance(*inst)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = launcher.ApplyMods(&allMods, utils.DWORD(shippingProcess.Pid), inst, launcher.MOD_DLLFILE)
		if err != nil {
			shippingProcess.Kill()
			launcherProcess.Kill()
			eacProcess.Kill()
			fmt.Println(err)
			return
		}

		shippingProcess.Wait()
		launcherProcess.Kill()
		eacProcess.Kill()
	}

	if err = utils.SaveData(&data); err != nil {
		fmt.Println(err)
		return
	}
}
