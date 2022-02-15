package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <demofile.dm2>\n", os.Args[0])
		return
	}

	demo := DemoFile{}
	demo.ParseDemo(os.Args[1])
	fmt.Printf("Map: %s (%s)\n", demo.Serverdata.MapName, demo.Configstrings[33].String)
	fmt.Printf("Frames: %d\n", len(demo.Frames))
}
