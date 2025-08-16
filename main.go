package main

import (
	"fmt"
	"os"
	"strings"

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
				allInstancesFormat = append(allInstancesFormat, fmt.Sprintf("	- %s (%s)", v.Name, v.BuildPath))
			}
			fmt.Printf("All instances (%d):\n%s", len(allInstancesFormat), strings.Join(allInstancesFormat, "\n"))
		}
	case "play":
		if len(opt.Parameters) != 1 {
			fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
			break
		}
		inst := utils.FetchInstance(&data, opt.Parameters[0])
		shipping, launcher, eac, err := utils.LaunchInstance(*inst)
		if err != nil {
			fmt.Println(err)
			break
		}

		shipping.Wait()
		launcher.Kill()
		eac.Kill()
	}

	if err = utils.SaveData(&data); err != nil {
		fmt.Println(err)
		return
	}
}
