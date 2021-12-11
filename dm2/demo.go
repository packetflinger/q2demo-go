package dm2

import (
	//"encoding/hex"

	"fmt"
	"os"

	"github.com/packetflinger/q2demo/msg"
)

const (
	MaxConfigStrings = 2080
	MaxEntities      = 1024
)

/**
 * A structure to store a parsed (packed) demo file
 */
type DemoFile struct {
	ParsingFrames bool
	Serverdata    msg.ServerData
	Configstrings [MaxConfigStrings]msg.ConfigString
	Baselines     [MaxEntities]msg.PackedEntity
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

	lenbuf := msg.MessageBuffer{Buffer: len, Index: 0}
	length := msg.ReadLong(&lenbuf)
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
	buf := msg.MessageBuffer{Buffer: lump}

	for buf.Index < len(buf.Buffer) {
		cmd := msg.ReadByte(&buf)

		switch cmd {
		case msg.SVCServerData:
			s := msg.ParseServerData(&buf)
			demo.Serverdata = s
			//fmt.Println(d)
			//fmt.Println(s)

		case msg.SVCConfigString:
			cs := msg.ParseConfigString(&buf)
			if !demo.ParsingFrames {
				demo.Configstrings[cs.Index] = cs
			}
			//fmt.Println(cs)

		case msg.SVCSpawnBaseline:
			_ = msg.ParseSpawnBaseline(&buf)
			//fmt.Println(bl)

		case msg.SVCStuffText:
			st := msg.ParseStuffText(&buf)
			// a "precache" stuff is the delimiter between header data
			// and frames
			if st.String == "precache\n" {
				demo.ParsingFrames = true
			}
			//fmt.Println(st)

		case msg.SVCFrame:
			_ = msg.ParseFrame(&buf)
			//fmt.Println(fr)

		case msg.SVCPlayerInfo:
			_ = msg.ParsePlayerstate(&buf)
			//fmt.Println(ps)

		case msg.SVCPacketEntities:
			_ = msg.ParsePacketEntities(&buf)
			//fmt.Println(ents)

		case msg.SVCPrint:
			_ = msg.ParsePrint(&buf)
			//fmt.Println(pr)
		}
	}
}

func ParseDemo(filename string, demo *DemoFile) {
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
