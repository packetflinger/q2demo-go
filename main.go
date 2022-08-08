package main

import (
	"flag"
)

type Flags struct {
	InputFile *string
	Verbose   *bool
	SVG       *bool
	Prints    *bool
	Layouts   *bool
	CStrings  *bool
}

var cli_args Flags

func main() {
	if *cli_args.InputFile == "" {
		flag.Usage()
		return
	}

	demo := DemoFile{}
	demo.ParseDemo(*cli_args.InputFile)
	//fmt.Printf("Map: %s (%s)\n", demo.Serverdata.MapName, demo.Configstrings[CSMapname].String)
	//fmt.Printf("Frames: %d\n", len(demo.Frames))

	//demo.WriteFile(demo.Filename + ".2")
	if *cli_args.SVG {
		demo.WriteIntermissionSVG()
	}
}

func init() {
	cli_args.InputFile = flag.String("i", "", "The input .dm2 file to work with")
	cli_args.SVG = flag.Bool("s", false, "Generate an SVG 'screenshot' of the intermission scoreboard")
	cli_args.Verbose = flag.Bool("v", false, "Show verbose output")
	cli_args.Prints = flag.Bool("p", false, "Output prints (console log)")
	cli_args.Layouts = flag.Bool("l", false, "Output layouts")
	cli_args.CStrings = flag.Bool("c", false, "Output Configstrings")
	flag.Parse()

	// don't double output prints
	if *cli_args.Verbose {
		*cli_args.Prints = false
		*cli_args.CStrings = false
	}
}
