package main

import (
	"fmt"
	"os"
)

var (
	playPath string
	allTweaks []string
	bIsServer bool = false
	redirectUrl string = ""
)

func main() {
	cobaltPath := "./redirector.dll"
	_, err := os.Stat(cobaltPath)
	if err != nil {
		fmt.Println("ERROR: Failed to find redirector DLL!")
		return
	}

	for i, v := range os.Args {
		if i == 0 {
			continue
		}

		if v == "--play" || v == "-p" {
			playPath = os.Args[i]
		} else if v == "--tweak" || v == "-t" {
			allTweaks = append(allTweaks, os.Args[i+1])
		} else if v == "--server" || v == "-s" {
			bIsServer = true
		} else if v == "--redirect" {
			redirectUrl = os.Args[i+1]
		}
	}

	if bIsServer {
		if _, err = os.Stat("./server.dll"); err != nil {
			fmt.Println("ERROR: Failed to find server DLL!")
			return
		}
	}


}
