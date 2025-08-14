package main

import (
	"fmt"
	"os"

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

	//TODO: replace with switch statement
	fmt.Println(opt)
	if opt.MainOperation == "profile" {
		if opt.SubOperation == "add" {
			if len(opt.Parameters) < 2 {
				fmt.Printf("Wrong parameter count! %d provided, 2 or more required!\n", len(opt.Parameters))
				return
			}
			utils.AddInstance(&data, opt.Parameters[0], opt.Parameters[1])
		}

		if opt.SubOperation == "remove" {
			if len(opt.Parameters) != 1 {
				fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
				return
			}
			utils.RemoveInstance(&data, opt.Parameters[0])
		}
	}

	if err = utils.SaveData(&data); err != nil {
		fmt.Println(err)
	}
}
