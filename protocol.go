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

func ReadLong(msg *MessageBuffer) int32 {
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

func WriteLong(data int32, msg *MessageBuffer) {
	msg.Buffer[msg.Index] = byte(data & 0xff)
	msg.Buffer[msg.Index+1] = byte((data >> 8) & 0xff)
	msg.Buffer[msg.Index+2] = byte((data >> 16) & 0xff)
	msg.Buffer[msg.Index+3] = byte((data >> 24) & 0xff)
	msg.Index += 4
}

/**
 * basically just grab a subsection of the buffer
 */
func ReadData(msg *MessageBuffer, length int) []byte {
	start := msg.Index
	msg.Index += length
	return msg.Buffer[start:msg.Index]
}

func WriteData(data []byte, msg *MessageBuffer) {
	msg.Buffer = append(msg.Buffer, data...)
	msg.Index += len(data)
}

/**
 * Keep building a string until we hit a null
 */
func ReadString(msg *MessageBuffer) string {
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

func WriteString(s string, msg *MessageBuffer) {
	b := []byte(s)
	msg.Buffer = append(msg.Buffer, b...)
	msg.Index += len(b)
}

/**
 * Read two bytes as a Short
 */
func ReadShort(msg *MessageBuffer) uint16 {
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

func WriteShort(s uint16, msg *MessageBuffer) {
	msg.Buffer[msg.Index] = byte(s & 0xff)
	msg.Buffer[msg.Index+1] = byte(s>>8) & 0xff
	msg.Index += 2
}

// for consistency
func ReadByte(msg *MessageBuffer) byte {
	val := byte(msg.Buffer[msg.Index])
	msg.Index++
	return val
}

func WriteByte(b byte, msg *MessageBuffer) {
	msg.Buffer[msg.Index] = b
	msg.Index++
}

func ReadChar(msg *MessageBuffer) int8 {
	val := int8(msg.Buffer[msg.Index])
	msg.Index++
	return val
}

func WriteChar(c uint8, msg *MessageBuffer) {
	msg.Buffer[msg.Index] = byte(c)
	msg.Index++
}

func ReadWord(msg *MessageBuffer) int16 {
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

func WriteWord(w int16, msg *MessageBuffer) {
	msg.Buffer[msg.Index] = byte(w & 0xff)
	msg.Buffer[msg.Index+1] = byte(w >> 8)
	msg.Index += 2
}
