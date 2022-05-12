package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	SVCBad = iota
	SVCMuzzleFlash
	SVCMuzzleFlash2
	SVCTempEntity
	SVCLayout
	SVCInventory
	SVCNOP
	SVCDisconnect
	SVCReconnect
	SVCSound
	SVCPrint
	SVCStuffText
	SVCServerData
	SVCConfigString
	SVCSpawnBaseline
	SVCCenterPrint
	SVCDownload
	SVCPlayerInfo
	SVCPacketEntities
	SVCDeltaPacketEntities
	SVCFrame
	SVCZPacket   // r1q2
	SVCZDownload // r1q2
	SVCGameState // r1q2/q2pro
	SVCSetting   // r1q2/q2pro
	SVCNumTypes  // r1q2/q2pro
)

// entity state flags
const (
	EntityOrigin1   = 1 << 0
	EntityOrigin2   = 1 << 1
	EntityAngle2    = 1 << 2
	EntityAngle3    = 1 << 3
	EntityFrame8    = 1 << 4
	EntityEvent     = 1 << 5
	EntityRemove    = 1 << 6
	EntityMoreBits1 = 1 << 7

	EntityNumber16  = 1 << 8
	EntityOrigin3   = 1 << 9
	EntityAngle1    = 1 << 10
	EntityModel     = 1 << 11
	EntityRenderFX8 = 1 << 12
	EntityAngle16   = 1 << 13
	EntityEffects8  = 1 << 14
	EntityMoreBits2 = 1 << 15

	EntitySkin8      = 1 << 16
	EntityFrame16    = 1 << 17
	EntityRenderFX16 = 1 << 18
	EntityEffects16  = 1 << 19
	EntityModel2     = 1 << 20
	EntityModel3     = 1 << 21
	EntityModel4     = 1 << 22
	EntityMoreBits3  = 1 << 23

	EntityOldOrigin = 1 << 24
	EntitySkin16    = 1 << 25
	EntitySound     = 1 << 26
	EntitySolid     = 1 << 27
)

const (
	PlayerType        = 1 << 0
	PlayerOrigin      = 1 << 1
	PlayerVelocity    = 1 << 2
	PlayerTime        = 1 << 3
	PlayerFlags       = 1 << 4
	PlayerGravity     = 1 << 5
	PlayerDeltaAngles = 1 << 6
	PlayerViewOffset  = 1 << 7

	PlayerViewAngles  = 1 << 8
	PlayerKickAngles  = 1 << 9
	PlayerBlend       = 1 << 10
	PlayerFOV         = 1 << 11
	PlayerWeaponIndex = 1 << 12
	PlayerWeaponFrame = 1 << 13
	PlayerRDFlags     = 1 << 14
	PlayerReserved    = 1 << 15

	PlayerBits = 16
	PlayerMask = (1 << PlayerBits) - 1
)

const (
	CSMapname = 33
)

func (msg *MessageBuffer) ReadLong() int32 {
	var tmp struct {
		Value int32
	}

	r := bytes.NewReader(msg.Buffer[msg.Index:])
	if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	msg.Index += 4
	return tmp.Value
}

func (msg *MessageBuffer) WriteLong(data int32) {
	msg.Buffer[msg.Index] = byte(data & 0xff)
	msg.Buffer[msg.Index+1] = byte((data >> 8) & 0xff)
	msg.Buffer[msg.Index+2] = byte((data >> 16) & 0xff)
	msg.Buffer[msg.Index+3] = byte((data >> 24) & 0xff)
	msg.Index += 4
}

/**
 * basically just grab a subsection of the buffer
 */
func (msg *MessageBuffer) ReadData(length int) []byte {
	start := msg.Index
	msg.Index += length
	return msg.Buffer[start:msg.Index]
}

func (msg *MessageBuffer) WriteData(data []byte) {
	msg.Buffer = append(msg.Buffer, data...)
	msg.Index += len(data)
}

/**
 * Keep building a string until we hit a null
 */
func (msg *MessageBuffer) ReadString() string {
	var buffer bytes.Buffer

	// find the next null (terminates the string)
	for i := 0; msg.Buffer[msg.Index] != 0; i++ {
		// we hit the end without finding a null
		if msg.Index == len(msg.Buffer) {
			break
		}

		buffer.WriteString(string(msg.Buffer[msg.Index]))
		msg.Index++
	}

	msg.Index++
	return buffer.String()
}

/**
 * Write a string char by char to a messagebuffer.
 * Strings are null terminated, so add a 0x00 at the end.
 */
func (msg *MessageBuffer) WriteString(s string) {
	for _, ch := range s {
		msg.WriteByte(byte(ch))
	}

	msg.WriteByte(0)
}

/**
 * Read two bytes as a Short
 */
func (msg *MessageBuffer) ReadShort() uint16 {
	var tmp struct {
		Value uint16
	}

	r := bytes.NewReader(msg.Buffer[msg.Index:])
	if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	msg.Index += 2
	return tmp.Value
}

func (msg *MessageBuffer) WriteShort(s uint16) {
	msg.Buffer[msg.Index] = byte(s & 0xff)
	msg.Buffer[msg.Index+1] = byte(s>>8) & 0xff
	msg.Index += 2
}

// for consistency
func (msg *MessageBuffer) ReadByte() byte {
	val := byte(msg.Buffer[msg.Index])
	msg.Index++
	return val
}

func (msg *MessageBuffer) WriteByte(b byte) {
	msg.Buffer[msg.Index] = b
	msg.Index++
}

func (msg *MessageBuffer) ReadChar() int8 {
	val := int8(msg.Buffer[msg.Index])
	msg.Index++
	return val
}

func (msg *MessageBuffer) WriteChar(c uint8) {
	msg.Buffer[msg.Index] = byte(c)
	msg.Index++
}

func (msg *MessageBuffer) ReadWord() int16 {
	var tmp struct {
		Value int16
	}

	r := bytes.NewReader(msg.Buffer[msg.Index:])
	if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	msg.Index += 2
	return tmp.Value
}

func (msg *MessageBuffer) WriteWord(w int16) {
	msg.Buffer[msg.Index] = byte(w & 0xff)
	msg.Buffer[msg.Index+1] = byte(w >> 8)
	msg.Index += 2
}
