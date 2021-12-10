package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/packetflinger/q2demo/dm2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <demofile.dm2>\n", os.Args[0])
		return
	}

	position := 0
	demofile := dm2.OpenDemo(os.Args[1])

	for {
		lump, size := dm2.NextLump(demofile, int64(position))
		if size == 0 {
			break
		}

		fmt.Printf("%s\n", hex.Dump(lump))
		position += size

		dm2.ParseLump(lump)
	}

	dm2.CloseDemo(demofile)
}
