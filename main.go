package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	InputFile *string
	Verbose   *bool
	SVG       *bool
}

var cli_args Flags

func main() {
	if *cli_args.InputFile == "" {
		flag.Usage()
		return
	}

	demo := DemoFile{}
	demo.ParseDemo(*cli_args.InputFile)
	fmt.Printf("Map: %s (%s)\n", demo.Serverdata.MapName, demo.Configstrings[CSMapname].String)
	fmt.Printf("Frames: %d\n", len(demo.Frames))

	//demo.WriteFile(demo.Filename + ".2")
}

func init() {
	cli_args.InputFile = flag.String("i", "", "The input .dm2 file to work with")
	cli_args.SVG = flag.Bool("s", false, "Generate an SVG 'screenshot' of the intermission scoreboard")
	cli_args.Verbose = flag.Bool("v", false, "Show verbose output")
	flag.Parse()
}
