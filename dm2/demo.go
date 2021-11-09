package dm2

import (
	//"encoding/hex"
	"fmt"
	"os"

	"github.com/packetflinger/q2demo/msg"
)

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

	fmt.Printf("Lump position: %d, length: %d\n", pos, length)
	_, err = f.Seek(pos+4, 0)
	check(err)

	lump := make([]byte, length)
	read, err := f.Read(lump)
	check(err)

	return lump, read+4
}

func ParseLump(lump []byte) {
	buf := msg.MessageBuffer{Buffer:lump}

	for buf.Index < len(buf.Buffer) {
		cmd := msg.ReadByte(&buf)

		switch (cmd) {
		case msg.SVCServerData:
			s := msg.ParseServerData(&buf)
			fmt.Println(s)

		case msg.SVCConfigString:
			cs := msg.ParseConfigString(&buf)
			fmt.Println(cs)

		case msg.SVCSpawnBaseline:
			bl := msg.ParseSpawnBaseline(&buf)
			fmt.Println(bl)

		case msg.SVCStuffText:
			st := msg.ParseStuffText(&buf)
			fmt.Println(st)
		}
	}
}