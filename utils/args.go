package utils

import (
	"fmt"
	"os"
)

type ArgOptions struct {
	MainOperation string
	SubOperation  string
	Parameters    []string
}

func ProcessArgs(args []string) (ArgOptions, error) {
	opts := ArgOptions{}

	if len(args) > 2 {
		opts.MainOperation = os.Args[1]
		if opts.MainOperation != "play" {
			opts.SubOperation = os.Args[2]
			opts.Parameters = args[3:]
		} else {
			opts.Parameters = args[2:]
		}
	} else {
		return ArgOptions{}, fmt.Errorf("Arguments must be at least a length of 2!")
	}

	return opts, nil
}
