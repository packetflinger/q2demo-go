package main

import (
	"encoding/hex"
	"fmt"

	"github.com/packetflinger/q2demo/dm2"
)

func main() {
	position := 0
	demofile := dm2.OpenDemo("test.dm2")

	for {
		lump, size := dm2.NextLump(demofile, int64(position))
		if size == 0 {
			break;
		}

		fmt.Printf("%s\n", hex.Dump(lump))
		position += size

		dm2.ParseLump(lump)
	}

	dm2.CloseDemo(demofile)
}
