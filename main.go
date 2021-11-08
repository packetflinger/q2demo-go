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
		lump := dm2.NextLump((demofile), int64(position))
		if len(lump) == 0 {
			break
		}

		position += len(lump)
		fmt.Printf("%s\n", hex.Dump(lump))
	}

	dm2.CloseDemo(demofile)
}
