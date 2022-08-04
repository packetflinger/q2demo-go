package main

import (
	//"encoding/hex"

	"encoding/hex"
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
	Filename      string
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
	length := lenbuf.ReadLong()
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
		cmd := buf.ReadByte()

		switch cmd {
		case SVCServerData:
			s := buf.ParseServerData()
			demo.Serverdata = s

		case SVCConfigString:
			cs := buf.ParseConfigString()
			if !demo.ParsingFrames {
				demo.Configstrings[cs.Index] = cs
			} else {
				currentframe.Strings = append(currentframe.Strings, cs)
			}

		case SVCSpawnBaseline:
			bl := buf.ParseSpawnBaseline()
			demo.Baselines[bl.Number] = bl

		case SVCStuffText:
			st := buf.ParseStuffText()
			// a "precache" stuff is the delimiter between header data
			// and frames
			if st.String == "precache\n" {
				demo.ParsingFrames = true
			}

		case SVCFrame:
			fr := buf.ParseFrame()
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
			ps := buf.ParseDeltaPlayerstate()
			currentframe.Playerstate = ps

		case SVCPacketEntities:
			ents := buf.ParsePacketEntities()
			for _, e := range ents {
				currentframe.Entities[e.Number] = e
			}

		case SVCPrint:
			_ = buf.ParsePrint()
		}
	}
}

func (demo *DemoFile) ParseDemo(filename string) {
	demo.Filename = filename
	position := 0
	demofile := OpenDemo(filename)

	for {
		lump, size := NextLump(demofile, int64(position))
		if size == 0 {
			break
		}

		position += size

		ParseLump(lump, demo)
	}

	CloseDemo(demofile)
}

/**
 * Build a valid .dm2 file from the demo structure
 */
func (demo *DemoFile) WriteFile(filename string) {
	// write the serverdata first
	msg := MessageBuffer{}
	msg.Buffer = make([]byte, 0xffff)
	msg.WriteByte(SVCServerData)
	msg.WriteLong(demo.Serverdata.Protocol)
	msg.WriteLong(demo.Serverdata.ServerCount)
	msg.WriteByte(1) // this is a demo
	msg.WriteString(demo.Serverdata.GameDir)
	msg.WriteShort(uint16(demo.Serverdata.ClientNumber))
	msg.WriteString(demo.Serverdata.MapName)

	// configstrings
	for _, cs := range demo.Configstrings {
		if cs.String == "" {
			continue
		}

		msg.WriteByte(SVCConfigString)
		msg.WriteShort(uint16(cs.Index))
		msg.WriteString(cs.String)
	}

	// baselines
	for _, ent := range demo.Baselines {
		if ent.Number == 0 {
			continue
		}

		msg.WriteByte(SVCSpawnBaseline)
		msg.WriteDeltaEntity(PackedEntity{}, ent)
	}

	msg.WriteByte(SVCStuffText)
	msg.WriteString("precache\n")

	previousframe = &ServerFrame{}
	previousframe.Frame.Number = -1 // no delta from first frame

	for _, fr := range demo.Frames {
		msg.WriteDeltaFrame(previousframe, &fr)
		previousframe = &fr
	}

	fmt.Printf("%s\n", hex.Dump(msg.Buffer[:msg.Index]))
}

// Parse the layout of the intermission screen and recreate it as an SVG
func (demo *DemoFile) WriteIntermissionSVG() {
	fmt.Println("Writing SVG!")
}
