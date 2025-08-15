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

	//TODO: Replace with switch statement

	switch opt.MainOperation {
	case "profile":
		switch opt.SubOperation {
		case "add":
			if len(opt.Parameters) < 2 {
				fmt.Printf("Wrong parameter count! %d provided, 2 or more required!\n", len(opt.Parameters))
				break // or return?
			}

			err := utils.CreateInstance(&data, opt.Parameters[0], opt.Parameters[1])
			if err != nil {
				fmt.Println(err)
				break // or return?
			}
		case "remove":
			if len(opt.Parameters) != 1 {
				fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
				break // or return?
			}
			utils.RemoveInstance(&data, opt.Parameters[0])
		}
	case "play":
		if len(opt.Parameters) != 1 {
			fmt.Printf("Wrong parameter count! %d provided, 1 required!\n", len(opt.Parameters))
			break // or return?
		}
		inst := utils.FetchInstance(&data, opt.Parameters[0])
		err := utils.LaunchInstance(*inst)
		if err != nil {
			fmt.Println(err)
			break // or return?
		}
	}

	if err = utils.SaveData(&data); err != nil {
		fmt.Println(err)
	}
}
