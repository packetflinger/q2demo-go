package dm2

import (
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

func NextLump(f *os.File, pos int64) []byte {
	_, err := f.Seek(pos, 0)
	check(err)

	len := make([]byte, 4)
	_, err = f.Read(len)
	check(err)

	lenbuf := msg.MessageBuffer{Buffer: len, Index: 0}
	length := msg.ReadLong(&lenbuf)
	if length == -1 {
		return []byte{}
	}

	_, err = f.Seek(pos+4, 0)
	check(err)

	lump := make([]byte, length)
	_, err = f.Read(lump)
	check(err)

	return lump
}
