package main

import (
	"fmt"
	"os"

	"github.com/packetflinger/q2demo/dm2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <demofile.dm2>\n", os.Args[0])
		return
	}

	demo := dm2.DemoFile{}
	dm2.ParseDemo(os.Args[1], &demo)
	fmt.Printf("Map: %s (%s)\n", demo.Serverdata.MapName, demo.Configstrings[33].String)
}
