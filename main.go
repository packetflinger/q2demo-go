package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	InputFile  *string
	Verbose    *bool
	SVG        *bool
	Prints     *bool
	Layouts    *bool
	CStrings   *bool
	Details    *bool
	OutputFile *string
}

var cli_args Flags

func main() {
	if *cli_args.InputFile == "" {
		flag.Usage()
		return
	}

	demo := DemoFile{}
	demo.ParseDemo(*cli_args.InputFile)

	if *cli_args.Details {
		framecount := len(demo.Frames)
		mins := int(framecount / 10 / 60)
		secs := int((framecount / 10) - (mins * 60))
		fmt.Printf("Map: %s (%s)\n", demo.Serverdata.MapName, demo.Configstrings[CSMapname].String)
		fmt.Printf("POV Entity: %d\n", demo.Serverdata.ClientNumber)
		fmt.Printf("Frames Count: %d\n", framecount)
		fmt.Printf("Length: %02d:%02d\n", mins, secs)
	}

	if *cli_args.SVG {
		demo.WriteIntermissionSVG()
	}

	// write the demo structure back to a new file
	if *cli_args.OutputFile != "" {
		demo.WriteFile(*cli_args.OutputFile)
	}
}

func init() {
	cli_args.InputFile = flag.String("i", "", "The input .dm2 file to work with")
	cli_args.SVG = flag.Bool("s", false, "Generate an SVG 'screenshot' of the intermission scoreboard")
	cli_args.Verbose = flag.Bool("v", false, "Show verbose output")
	cli_args.Prints = flag.Bool("p", false, "Output prints (console log)")
	cli_args.Layouts = flag.Bool("l", false, "Output layouts")
	cli_args.CStrings = flag.Bool("c", false, "Output Configstrings")
	cli_args.Details = flag.Bool("d", false, "Show details about the demo after parsing")
	cli_args.OutputFile = flag.String("o", "", "The output .dm2 file")
	flag.Parse()

	// manually change conflicting flags
	if *cli_args.Verbose {
		*cli_args.Prints = false
		*cli_args.CStrings = false
		*cli_args.Layouts = false
	}
}
