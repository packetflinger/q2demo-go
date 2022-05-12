package main

type MessageBuffer struct {
	Buffer []byte
	Index  int
	Length int // maybe not needed
}

type ServerData struct {
	Protocol     int32
	ServerCount  int32
	Demo         int8
	GameDir      string
	ClientNumber int16
	MapName      string
}

type ConfigString struct {
	Index  int16
	String string
}

type StuffText struct {
	String string
}

type PackedEntity struct {
	Number      uint32
	Origin      [3]int16
	Angles      [3]int16
	OldOrigin   [3]int16
	ModelIndex  uint8
	ModelIndex2 uint8
	ModelIndex3 uint8
	ModelIndex4 uint8
	SkinNum     uint32
	Effects     uint32
	RenderFX    uint32
	Solid       uint32
	Frame       uint16
	Sound       uint8
	Event       uint8
}

type PlayerMoveState struct {
	Type        uint8
	Origin      [3]int16
	Velocity    [3]int16
	Flags       byte
	Time        byte
	Gravity     int16
	DeltaAngles [3]int16
}

type PackedPlayer struct {
	PlayerMove PlayerMoveState
	ViewAngles [3]int16
	ViewOffset [3]int8
	KickAngles [3]int8
	GunAngles  [3]int8
	GunOffset  [3]int8
	GunIndex   uint8
	GunFrame   uint8
	Blend      [4]uint8
	FOV        uint8
	RDFlags    uint8
	Stats      [32]int16
}

type FrameMsg struct {
	Number     int32
	Delta      int32
	Suppressed int8
	AreaBytes  int8
	AreaBits   []byte
}

type Print struct {
	Level  uint8
	String string
}

func (m *MessageBuffer) ParseServerData() ServerData {
	sd := ServerData{}

	sd.Protocol = m.ReadLong()
	sd.ServerCount = m.ReadLong()
	sd.Demo = int8(m.ReadByte())
	sd.GameDir = m.ReadString()
	sd.ClientNumber = int16(m.ReadShort())
	sd.MapName = m.ReadString()

	return sd
}

func (m *MessageBuffer) ParseConfigString() ConfigString {
	cs := ConfigString{
		Index:  int16(m.ReadShort()),
		String: m.ReadString(),
	}

	return cs
}

func ParseSpawnBaseline(m *MessageBuffer) PackedEntity {
	bitmask := ParseEntityBitmask(m)
	number := ParseEntityNumber(m, bitmask)
	ent := ParseEntity(m, PackedEntity{}, number, bitmask)

	return ent
}

func ParseEntityBitmask(m *MessageBuffer) uint32 {
	bits := uint32(m.ReadByte())

	if bits&EntityMoreBits1 != 0 {
		bits |= (uint32(m.ReadByte()) << 8)
	}

	if bits&EntityMoreBits2 != 0 {
		bits |= (uint32(m.ReadByte()) << 16)
	}

	if bits&EntityMoreBits3 != 0 {
		bits |= (uint32(m.ReadByte()) << 24)
	}

	return uint32(bits)
}

func ParseEntityNumber(m *MessageBuffer, flags uint32) uint16 {
	num := uint16(0)
	if flags&EntityNumber16 != 0 {
		num = uint16(m.ReadShort())
	} else {
		num = uint16(m.ReadByte())
	}

	return num
}

func ParseEntity(m *MessageBuffer, from PackedEntity, num uint16, bits uint32) PackedEntity {
	to := from
	to.Number = uint32(num)

	if bits == 0 {
		return to
	}

	if bits&EntityModel != 0 {
		to.ModelIndex = uint8(m.ReadByte())
	}

	if bits&EntityModel2 != 0 {
		to.ModelIndex2 = uint8(m.ReadByte())
	}

	if bits&EntityModel3 != 0 {
		to.ModelIndex3 = uint8(m.ReadByte())
	}

	if bits&EntityModel4 != 0 {
		to.ModelIndex4 = uint8(m.ReadByte())
	}

	if bits&EntityFrame8 != 0 {
		to.Frame = uint16(m.ReadByte())
	}

	if bits&EntityFrame16 != 0 {
		to.Frame = uint16(m.ReadShort())
	}

	if (bits & (EntitySkin8 | EntitySkin16)) == (EntitySkin8 | EntitySkin16) {
		to.SkinNum = uint32(m.ReadLong())
	} else if bits&EntitySkin8 != 0 {
		to.SkinNum = uint32(m.ReadByte())
	} else if bits&EntitySkin16 != 0 {
		to.SkinNum = uint32(m.ReadWord())
	}

	if (bits & (EntityEffects8 | EntityEffects16)) == (EntityEffects8 | EntityEffects16) {
		to.Effects = uint32(m.ReadLong())
	} else if bits&EntityEffects8 != 0 {
		to.Effects = uint32(m.ReadByte())
	} else if bits&EntityEffects16 != 0 {
		to.Effects = uint32(m.ReadWord())
	}

	if (bits & (EntityRenderFX8 | EntityRenderFX16)) == (EntityRenderFX8 | EntityRenderFX16) {
		to.RenderFX = uint32(m.ReadLong())
	} else if bits&EntityRenderFX8 != 0 {
		to.RenderFX = uint32(m.ReadByte())
	} else if bits&EntityRenderFX16 != 0 {
		to.RenderFX = uint32(m.ReadWord())
	}

	if bits&EntityOrigin1 != 0 {
		to.Origin[0] = int16(m.ReadShort())
	}

	if bits&EntityOrigin2 != 0 {
		to.Origin[1] = int16(m.ReadShort())
	}

	if bits&EntityOrigin3 != 0 {
		to.Origin[2] = int16(m.ReadShort())
	}

	if bits&EntityAngle1 != 0 {
		to.Angles[0] = int16(m.ReadByte())
	}

	if bits&EntityAngle2 != 0 {
		to.Angles[1] = int16(m.ReadByte())
	}

	if bits&EntityAngle3 != 0 {
		to.Angles[2] = int16(m.ReadByte())
	}

	if bits&EntityOldOrigin != 0 {
		to.OldOrigin[0] = int16(m.ReadShort())
		to.OldOrigin[1] = int16(m.ReadShort())
		to.OldOrigin[2] = int16(m.ReadShort())
	}

	if bits&EntitySound != 0 {
		to.Sound = uint8(m.ReadByte())
	}

	if bits&EntityEvent != 0 {
		to.Event = uint8(m.ReadByte())
	}

	if bits&EntitySolid != 0 {
		to.Solid = uint32(m.ReadWord())
	}

	return to
}

func ParseStuffText(m *MessageBuffer) StuffText {
	str := StuffText{String: m.ReadString()}
	return str
}

func ParseFrame(m *MessageBuffer) FrameMsg {
	N := int32(m.ReadLong())
	D := int32(m.ReadLong())
	S := int8(m.ReadByte())
	A := int8(m.ReadByte())
	Ab := m.ReadData(int(A))

	fr := FrameMsg{
		Number:     N,
		Delta:      D,
		Suppressed: S,
		AreaBytes:  A,
		AreaBits:   Ab,
	}

	return fr
}

func ParseDeltaPlayerstate(m *MessageBuffer) PackedPlayer {
	bits := m.ReadWord()
	pm := PlayerMoveState{}
	ps := PackedPlayer{}

	/*if lastframe != nil {
		ps = lastframe.Playerstate
	}*/

	if bits&PlayerType != 0 {
		pm.Type = uint8(m.ReadByte())
	}

	if bits&PlayerOrigin != 0 {
		pm.Origin[0] = int16(m.ReadShort())
		pm.Origin[1] = int16(m.ReadShort())
		pm.Origin[2] = int16(m.ReadShort())
	}

	if bits&PlayerVelocity != 0 {
		pm.Velocity[0] = int16(m.ReadShort())
		pm.Velocity[1] = int16(m.ReadShort())
		pm.Velocity[2] = int16(m.ReadShort())
	}

	if bits&PlayerTime != 0 {
		pm.Time = m.ReadByte()
	}

	if bits&PlayerFlags != 0 {
		pm.Flags = m.ReadByte()
	}

	if bits&PlayerGravity != 0 {
		pm.Gravity = int16(m.ReadShort())
	}

	if bits&PlayerDeltaAngles != 0 {
		pm.DeltaAngles[0] = int16(m.ReadShort())
		pm.DeltaAngles[1] = int16(m.ReadShort())
		pm.DeltaAngles[2] = int16(m.ReadShort())
	}

	if bits&PlayerViewOffset != 0 {
		ps.ViewOffset[0] = int8(m.ReadChar())
		ps.ViewOffset[1] = int8(m.ReadChar())
		ps.ViewOffset[2] = int8(m.ReadChar())
	}

	if bits&PlayerViewAngles != 0 {
		ps.ViewAngles[0] = int16(m.ReadShort())
		ps.ViewAngles[1] = int16(m.ReadShort())
		ps.ViewAngles[2] = int16(m.ReadShort())
	}

	if bits&PlayerWeaponIndex != 0 {
		ps.GunIndex = uint8(m.ReadByte())
	}

	if bits&PlayerWeaponFrame != 0 {
		ps.GunFrame = uint8(m.ReadByte())
		ps.GunOffset[0] = int8(m.ReadChar())
		ps.GunOffset[1] = int8(m.ReadChar())
		ps.GunOffset[2] = int8(m.ReadChar())
		ps.GunAngles[0] = int8(m.ReadChar())
		ps.GunAngles[1] = int8(m.ReadChar())
		ps.GunAngles[2] = int8(m.ReadChar())
	}

	if bits&PlayerBlend != 0 {
		ps.Blend[0] = uint8(m.ReadChar())
		ps.Blend[1] = uint8(m.ReadChar())
		ps.Blend[2] = uint8(m.ReadChar())
		ps.Blend[3] = uint8(m.ReadChar())
	}

	if bits&PlayerFOV != 0 {
		ps.FOV = uint8(m.ReadByte())
	}

	if bits&PlayerRDFlags != 0 {
		ps.RDFlags = uint8(m.ReadByte())
	}

	statbits := int32(m.ReadLong())
	for i := 0; i < 32; i++ {
		if statbits&(1<<i) != 0 {
			ps.Stats[i] = int16(m.ReadShort())
		}
	}

	ps.PlayerMove = pm

	return ps
}

func ParsePacketEntities(m *MessageBuffer) []PackedEntity {
	ents := []PackedEntity{}
	for {
		bits := ParseEntityBitmask(m)
		num := ParseEntityNumber(m, bits)

		if num <= 0 {
			break
		}

		entity := ParseEntity(m, PackedEntity{}, num, bits)
		ents = append(ents, entity)
	}

	return ents
}

func ParsePrint(m *MessageBuffer) Print {
	st := Print{
		Level:  uint8(m.ReadByte()),
		String: m.ReadString(),
	}

	return st
}
