package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	switch opt.MainOperation {
	case "profiles":
		switch opt.SubOperation {
		case "add":
			if len(opt.Parameters) < 2 {
				fmt.Printf("Wrong parameter count! %d provided, 2 or more required!\n", len(opt.Parameters))
				break
			}

			err := utils.CreateInstance(&data.Instances, opt.Parameters[0], opt.Parameters[1])
			if err != nil {
				fmt.Println(err)
				break
			}
		case "remove":
			if len(opt.Parameters) != 1 {
				fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
				break
			}
			utils.RemoveInstance(&data.Instances, opt.Parameters[0])
		case "list":
			allInstancesFormat := []string{}
			for _, v := range data.Instances {
				allInstancesFormat = append(allInstancesFormat, fmt.Sprintf("\t- %s (%s)", v.Name, v.BuildPath))
			}
			fmt.Printf("All instances (%d):\n%s", len(allInstancesFormat), strings.Join(allInstancesFormat, "\n"))
		case "select":
			if len(opt.Parameters) != 1 {
				fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
				break
			}

			inst := utils.FetchInstance(&data.Instances, opt.Parameters[0])
			if inst == nil {
				fmt.Printf("Failed to fetch instance of name %s\n", opt.Parameters[0])
				break
			}
			data.SelectedInstance = inst.Name
		}
	case "play":
		if len(opt.Parameters) != 1 {
			fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
			break
		}

		inst := utils.FetchInstance(&data.Instances, opt.Parameters[0])
		allMods, err := launcher.CollectMods(inst.ModsPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = launcher.ApplyPaks(&allMods, inst)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Wait for patches later
		// _, err = launcher.ApplyPatches(&allMods, inst)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		shippingProcess, launcherProcess, eacProcess, err := launcher.LaunchInstance(*inst)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = launcher.ApplyDLLs(&allMods, inst, utils.DWORD(shippingProcess.Pid))
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

		err = launcher.RemovePaks(&allMods, inst)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "mods":
		switch opt.SubOperation {
		case "add":
			if len(opt.Parameters) != 1 {
				fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
				break
			}

			if data.SelectedInstance == "" {
				fmt.Printf("No profile has been selected! You must select a profile before adding mods!\n")
				break
			}

			inst := utils.FetchInstance(&data.Instances, data.SelectedInstance)
			if inst == nil {
				fmt.Printf("Failed to find the selected instance!\n")
				break
			}

			modBase := filepath.Base(opt.Parameters[0])
			err = launcher.CopyFileContents(opt.Parameters[0], filepath.Join(inst.BuildPath, "Mods", modBase))
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Printf("Successfully copied %s!\n", modBase)
		}
	}

	if err = utils.SaveData(&data); err != nil {
		fmt.Println(err)
		return
	}
}
