package main

import (
	//"encoding/hex"

	"fmt"
	"os"
)

var currentframe *ServerFrame
var previousframe *ServerFrame

const (
	MaxConfigStrings = 2080
	MaxEntities      = 1024
)

type ServerFrame struct {
	Frame       FrameMsg
	Playerstate PackedPlayer
	Entities    [MaxEntities]PackedEntity
	Strings     []ConfigString
	Prints      []Print
	Stuffs      []StuffText
}

/**
 * A structure to store a parsed (packed) demo file
 */
type DemoFile struct {
	ParsingFrames bool
	Serverdata    ServerData
	Configstrings [MaxConfigStrings]ConfigString
	Baselines     [MaxEntities]PackedEntity
	Frames        []ServerFrame
}

var Demo DemoFile
var CurrentPosition int64 = 0

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func OpenDemo(filename string) *os.File {
	f, err := os.Open(filename)
	check(err)

	fmt.Println("Demo file:", filename)
	return f
}

func CloseDemo(f *os.File) {
	f.Close()
}

func NextLump(f *os.File, pos int64) ([]byte, int) {
	_, err := f.Seek(pos, 0)
	check(err)

	len := make([]byte, 4)
	_, err = f.Read(len)
	check(err)

	//fmt.Printf("%s\n", hex.Dump(len))

	lenbuf := MessageBuffer{Buffer: len, Index: 0}
	length := ReadLong(&lenbuf)
	if length == -1 {
		return []byte{}, 0
	}

	//fmt.Printf("Lump position: %d, length: %d\n", pos, length)
	_, err = f.Seek(pos+4, 0)
	check(err)

	lump := make([]byte, length)
	read, err := f.Read(lump)
	check(err)

	return lump, read + 4
}

func ParseLump(lump []byte, demo *DemoFile) {
	buf := MessageBuffer{Buffer: lump}

	for buf.Index < len(buf.Buffer) {
		cmd := ReadByte(&buf)

		switch cmd {
		case SVCServerData:
			s := ParseServerData(&buf)
			demo.Serverdata = s

		case SVCConfigString:
			cs := ParseConfigString(&buf)
			if !demo.ParsingFrames {
				demo.Configstrings[cs.Index] = cs
			} else {
				currentframe.Strings = append(currentframe.Strings, cs)
			}

		case SVCSpawnBaseline:
			bl := ParseSpawnBaseline(&buf)
			demo.Baselines[bl.Number] = bl

		case SVCStuffText:
			st := ParseStuffText(&buf)
			// a "precache" stuff is the delimiter between header data
			// and frames
			if st.String == "precache\n" {
				demo.ParsingFrames = true
			}

		case SVCFrame:
			fr := ParseFrame(&buf)
			demo.Frames = append(demo.Frames, ServerFrame{})
			if currentframe != nil {
				previousframe = currentframe
			}
			currentframe = &demo.Frames[len(demo.Frames)-1]
			currentframe.Frame = fr
			if previousframe != nil {
				currentframe.Playerstate = previousframe.Playerstate
				currentframe.Entities = previousframe.Entities
			}

		case SVCPlayerInfo:
			ps := ParseDeltaPlayerstate(&buf)
			currentframe.Playerstate = ps

		case SVCPacketEntities:
			ents := ParsePacketEntities(&buf)
			for _, e := range ents {
				currentframe.Entities[e.Number] = e
			}

		case SVCPrint:
			_ = ParsePrint(&buf)
		}
	}
}

func (demo *DemoFile) ParseDemo(filename string) {
	position := 0
	demofile := OpenDemo(filename)

	for {
		lump, size := NextLump(demofile, int64(position))
		if size == 0 {
			break
		}

		//fmt.Printf("%s\n", hex.Dump(lump))
		position += size

		ParseLump(lump, demo)
	}

	CloseDemo(demofile)
}

/**
 * Build a valid .dm2 file from the demo structure
 */
func CreateDemoFile(demo *DemoFile, filename string) {

}
