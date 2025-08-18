package utils

import (
	"fmt"
)

type ArgOptions struct {
	MainOperation string
	SubOperation  string
	Parameters    []string
}

func ProcessArgs(args []string) (ArgOptions, error) {
	opts := ArgOptions{}

	if len(args) > 2 {
		opts.MainOperation = args[1]
		if opts.MainOperation == "play" || opts.MainOperation == "help" {
			opts.Parameters = args[2:]
		} else {
			opts.SubOperation = args[2]
			opts.Parameters = args[3:]
		}
	} else {
		return ArgOptions{}, fmt.Errorf("arguments must be at least a length of 2")
	}

	return opts, nil
}
