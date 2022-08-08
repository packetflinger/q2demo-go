package main

import (
	"fmt"
)

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

type PackedSound struct {
	Flags       uint8
	Index       uint8
	Volume      uint8
	Attenuation uint8
	TimeOffset  uint8
	Channel     uint16
	Entity      uint16
	Position    [3]uint16
}

type TemporaryEntity struct {
	Type      uint8
	Position1 [3]uint16
	Position2 [3]uint16
	Offset    [3]uint16
	Direction uint8
	Count     uint8
	Color     uint8
	Entity1   int16
	Entity2   int16
	Time      int32
}

type MuzzleFlash struct {
	Entity uint16
	Weapon uint8
}

type Layout struct {
	Data string
}

type CenterPrint struct {
	Data string
}

func (m *MessageBuffer) ParseServerData() ServerData {
	sd := ServerData{}

	sd.Protocol = m.ReadLong()
	sd.ServerCount = m.ReadLong()
	sd.Demo = int8(m.ReadByte())
	sd.GameDir = m.ReadString()
	sd.ClientNumber = int16(m.ReadShort())
	sd.MapName = m.ReadString()

	if *cli_args.Verbose {
		fmt.Printf(" * ServerData [%d] %s\n", sd.ClientNumber, sd.MapName)
	}

	return sd
}

func (m *MessageBuffer) ParseConfigString() ConfigString {
	cs := ConfigString{
		Index:  int16(m.ReadShort()),
		String: m.ReadString(),
	}

	if *cli_args.Verbose {
		fmt.Printf(" * ConfigString [%d] %s\n", cs.Index, cs.String)
	}

	return cs
}

func (m *MessageBuffer) ParseSpawnBaseline() PackedEntity {
	bitmask := m.ParseEntityBitmask()
	number := m.ParseEntityNumber(bitmask)
	ent := m.ParseEntity(PackedEntity{}, number, bitmask)

	if *cli_args.Verbose {
		fmt.Printf(" * Baseline [%d]\n", number)
	}
	return ent
}

func (m *MessageBuffer) ParseEntityBitmask() uint32 {
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

func (m *MessageBuffer) ParseEntityNumber(flags uint32) uint16 {
	num := uint16(0)
	if flags&EntityNumber16 != 0 {
		num = uint16(m.ReadShort())
	} else {
		num = uint16(m.ReadByte())
	}

	return num
}

func (m *MessageBuffer) ParseEntity(from PackedEntity, num uint16, bits uint32) PackedEntity {
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

func (m *MessageBuffer) ParseStuffText() StuffText {
	str := StuffText{String: m.ReadString()}
	if *cli_args.Verbose {
		fmt.Printf(" * Stuff \"%s\"\n", str)
	}
	return str
}

func (m *MessageBuffer) ParseFrame() FrameMsg {
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

	if *cli_args.Verbose {
		fmt.Printf(" * Frame [%d,%d]\n", N, D)
	}
	return fr
}

func (m *MessageBuffer) ParseDeltaPlayerstate() PackedPlayer {
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

	if bits&PlayerKickAngles != 0 {
		ps.KickAngles[0] = int8(m.ReadChar())
		ps.KickAngles[1] = int8(m.ReadChar())
		ps.KickAngles[2] = int8(m.ReadChar())
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

	if *cli_args.Verbose {
		fmt.Printf(" * PlayerState [%d]\n", bits)
	}

	return ps
}

func (m *MessageBuffer) ParsePacketEntities() []PackedEntity {
	ents := []PackedEntity{}
	if *cli_args.Verbose {
		fmt.Printf(" * Entities\n")
	}

	for {
		bits := m.ParseEntityBitmask()
		num := m.ParseEntityNumber(bits)

		if num <= 0 {
			break
		}

		entity := m.ParseEntity(PackedEntity{}, num, bits)
		ents = append(ents, entity)

		if *cli_args.Verbose {
			fmt.Printf("   - Ent [%d]\n", entity.Number)
		}
	}

	return ents
}

func (m *MessageBuffer) ParsePrint() Print {
	st := Print{
		Level:  uint8(m.ReadByte()),
		String: m.ReadString(),
	}

	if *cli_args.Verbose {
		fmt.Printf(" * Print \"%s\"\n", st.String[:len(st.String)-1])
	}

	if *cli_args.Prints {
		fmt.Printf("%s\n", StripConsoleChars(st.String[:len(st.String)-1]))
	}
	return st
}

func (m *MessageBuffer) ParseSound() PackedSound {
	s := PackedSound{}
	s.Flags = m.ReadByte()
	s.Index = m.ReadByte()

	if (s.Flags & SoundVolume) > 0 {
		s.Volume = m.ReadByte()
	} else {
		s.Volume = 1
	}

	if (s.Flags & SoundAttenuation) > 0 {
		s.Attenuation = m.ReadByte()
	} else {
		s.Attenuation = 1
	}

	if (s.Flags & SoundOffset) > 0 {
		s.TimeOffset = m.ReadByte()
	} else {
		s.TimeOffset = 0
	}

	if (s.Flags & SoundEntity) > 0 {
		s.Channel = m.ReadShort() & 7
		s.Entity = s.Channel >> 3
	} else {
		s.Channel = 0
		s.Entity = 0
	}

	if (s.Flags & SoundPosition) > 0 {
		s.Position = m.ReadPosition()
	}

	if *cli_args.Verbose {
		fmt.Printf(" * Sound [%d/%d]\n", s.Entity, s.Volume)
	}

	return s
}

/**
 * Find the differences between these two Entities
 */
func (to *PackedEntity) DeltaEntityBitmask(from *PackedEntity) int {
	bits := 0
	mask := uint32(0xffff8000)

	if to.Origin[0] != from.Origin[0] {
		bits |= EntityOrigin1
	}

	if to.Origin[1] != from.Origin[1] {
		bits |= EntityOrigin2
	}

	if to.Origin[2] != from.Origin[2] {
		bits |= EntityOrigin3
	}

	if to.Angles[0] != from.Angles[0] {
		bits |= EntityAngle1
	}

	if to.Angles[1] != from.Angles[1] {
		bits |= EntityAngle2
	}

	if to.Angles[2] != from.Angles[2] {
		bits |= EntityAngle3
	}

	if to.SkinNum != from.SkinNum {
		if to.SkinNum&mask&mask > 0 {
			bits |= EntitySkin8 | EntitySkin16
		} else if to.SkinNum&uint32(0x0000ff00) > 0 {
			bits |= EntitySkin16
		} else {
			bits |= EntitySkin8
		}
	}

	if to.Frame != from.Frame {
		if to.Frame&uint16(0xff00) > 0 {
			bits |= EntityFrame16
		} else {
			bits |= EntityFrame8
		}
	}

	if to.Effects != from.Effects {
		if to.Effects&mask > 0 {
			bits |= EntityEffects8 | EntityEffects16
		} else if to.Effects&0x0000ff00 > 0 {
			bits |= EntityEffects16
		} else {
			bits |= EntityEffects8
		}
	}

	if to.RenderFX != from.RenderFX {
		if to.RenderFX&mask > 0 {
			bits |= EntityRenderFX8 | EntityRenderFX16
		} else if to.RenderFX&0x0000ff00 > 0 {
			bits |= EntityRenderFX16
		} else {
			bits |= EntityRenderFX8
		}
	}

	if to.Solid != from.Solid {
		bits |= EntitySolid
	}

	if to.Event != from.Event {
		bits |= EntityEvent
	}

	if to.ModelIndex != from.ModelIndex {
		bits |= EntityModel
	}

	if to.ModelIndex2 != from.ModelIndex2 {
		bits |= EntityModel2
	}

	if to.ModelIndex3 != from.ModelIndex3 {
		bits |= EntityModel3
	}

	if to.ModelIndex4 != from.ModelIndex4 {
		bits |= EntityModel4
	}

	if to.Sound != from.Sound {
		bits |= EntitySound
	}

	if to.RenderFX&RFFrameLerp > 0 {
		bits |= EntityOldOrigin
	} else if to.RenderFX&RFBeam > 0 {
		bits |= EntityOldOrigin
	}

	if to.Number&0xff00 > 0 {
		bits |= EntityNumber16
	}

	if bits&0xff000000 > 0 {
		bits |= EntityMoreBits3 | EntityMoreBits2 | EntityMoreBits1
	} else if bits&0x00ff0000 > 0 {
		bits |= EntityMoreBits2 | EntityMoreBits1
	} else if bits&0x0000ff00 > 0 {
		bits |= EntityMoreBits1
	}

	return bits
}

/**
 * Compare from and to and only write what's different.
 * This is "delta compression"
 */
func (m *MessageBuffer) WriteDeltaEntity(from PackedEntity, to PackedEntity) {
	bits := to.DeltaEntityBitmask(&from)

	// write the bitmask first
	m.WriteByte(byte(bits & 255))
	if bits&0xff000000 > 0 {
		m.WriteByte(byte((bits >> 8) & 255))
		m.WriteByte(byte((bits >> 16) & 255))
		m.WriteByte(byte((bits >> 24) & 255))
	} else if bits&0x00ff0000 > 0 {
		m.WriteByte(byte((bits >> 8) & 255))
		m.WriteByte(byte((bits >> 16) & 255))
	} else if bits&0x0000ff00 > 0 {
		m.WriteByte(byte((bits >> 8) & 255))
	}

	// write the edict number
	if bits&EntityNumber16 > 0 {
		m.WriteShort(uint16(to.Number))
	} else {
		m.WriteByte(byte(to.Number))
	}

	if bits&EntityModel > 0 {
		m.WriteByte(to.ModelIndex)
	}

	if bits&EntityModel2 > 0 {
		m.WriteByte(to.ModelIndex2)
	}

	if bits&EntityModel3 > 0 {
		m.WriteByte(to.ModelIndex3)
	}

	if bits&EntityModel4 > 0 {
		m.WriteByte(to.ModelIndex4)
	}

	if bits&EntityFrame8 > 0 {
		m.WriteByte(byte(to.Frame))
	} else if bits&EntityFrame16 > 0 {
		m.WriteShort(to.Frame)
	}

	if (bits & (EntitySkin8 | EntitySkin16)) == (EntitySkin8 | EntitySkin16) {
		m.WriteLong(int32(to.SkinNum))
	} else if bits&EntitySkin8 > 0 {
		m.WriteByte(byte(to.SkinNum))
	} else if bits&EntitySkin16 > 0 {
		m.WriteShort(uint16(to.SkinNum))
	}

	if (bits & (EntityEffects8 | EntityEffects16)) == (EntityEffects8 | EntityEffects16) {
		m.WriteLong(int32(to.Effects))
	} else if bits&EntityEffects8 > 0 {
		m.WriteByte(byte(to.Effects))
	} else if bits&EntityEffects16 > 0 {
		m.WriteShort(uint16(to.Effects))
	}

	if (bits & (EntityRenderFX8 | EntityRenderFX16)) == (EntityRenderFX8 | EntityRenderFX16) {
		m.WriteLong(int32(to.RenderFX))
	} else if bits&EntityRenderFX8 > 0 {
		m.WriteByte(byte(to.RenderFX))
	} else if bits&EntityRenderFX16 > 0 {
		m.WriteShort(uint16(to.RenderFX))
	}

	if bits&EntityOrigin1 > 0 {
		m.WriteShort(uint16(to.Origin[0]))
	}

	if bits&EntityOrigin2 > 0 {
		m.WriteShort(uint16(to.Origin[1]))
	}

	if bits&EntityOrigin3 > 0 {
		m.WriteShort(uint16(to.Origin[2]))
	}

	if bits&EntityAngle1 > 0 {
		m.WriteByte(byte(to.Angles[0] >> 8))
	}

	if bits&EntityAngle2 > 0 {
		m.WriteByte(byte(to.Angles[1] >> 8))
	}

	if bits&EntityAngle3 > 0 {
		m.WriteByte(byte(to.Angles[2] >> 8))
	}

	if bits&EntityOldOrigin > 0 {
		m.WriteShort(uint16(to.OldOrigin[0]))
		m.WriteShort(uint16(to.OldOrigin[1]))
		m.WriteShort(uint16(to.OldOrigin[2]))
	}

	if bits&EntitySound > 0 {
		m.WriteByte(to.Sound)
	}

	if bits&EntityEvent > 0 {
		m.WriteByte(to.Event)
	}

	if bits&EntitySolid > 0 {
		m.WriteShort(uint16(to.Solid))
	}
}

/**
 * compress frames
 */
func (m *MessageBuffer) WriteDeltaFrame(from *ServerFrame, to *ServerFrame) {
	m.WriteByte(SVCFrame)
	m.WriteLong(to.Frame.Number)
	m.WriteLong(from.Frame.Number)
	m.WriteByte(byte(to.Frame.Suppressed))
	m.WriteByte(byte(to.Frame.AreaBytes))
	m.WriteData(to.Frame.AreaBits)

}

func (m *MessageBuffer) ParseTempEntity() TemporaryEntity {
	te := TemporaryEntity{}

	te.Type = m.ReadByte()
	switch te.Type {
	case TentBlood:
		fallthrough
	case TentGunshot:
		fallthrough
	case TentSparks:
		fallthrough
	case TentBulletSparks:
		fallthrough
	case TentScreenSparks:
		fallthrough
	case TentShieldSparks:
		fallthrough
	case TentShotgun:
		fallthrough
	case TentBlaster:
		fallthrough
	case TentGreenBlood:
		fallthrough
	case TentBlaster2:
		fallthrough
	case TentFlechette:
		fallthrough
	case TentHeatBeamSparks:
		fallthrough
	case TentHeatBeamSteam:
		fallthrough
	case TentMoreBlood:
		fallthrough
	case TentElectricSparks:
		te.Position1 = m.ReadPosition()
		te.Direction = m.ReadDirection()
	case TentSplash:
		fallthrough
	case TentLaserSparks:
		fallthrough
	case TentWeldingSparks:
		fallthrough
	case TentTunnelSparks:
		te.Count = m.ReadByte()
		te.Position1 = m.ReadPosition()
		te.Direction = m.ReadDirection()
		te.Color = m.ReadByte()
	case TentBlueHyperBlaster:
		fallthrough
	case TentRailTrail:
		fallthrough
	case TentBubbleTrail:
		fallthrough
	case TentDebugTrail:
		fallthrough
	case TentBubbleTrail2:
		fallthrough
	case TentBFGLaser:
		te.Position1 = m.ReadPosition()
		te.Position2 = m.ReadPosition()
	case TentGrenadeExplosion:
		fallthrough
	case TentGrenadeExplosionWater:
		fallthrough
	case TentExplosion2:
		fallthrough
	case TentPlasmaExplosion:
		fallthrough
	case TentRocketExplosion:
		fallthrough
	case TentRocketExplosionWater:
		fallthrough
	case TentExplosion1:
		fallthrough
	case TentExplosion1NP:
		fallthrough
	case TentExplosion1Big:
		fallthrough
	case TentBFGExplosion:
		fallthrough
	case TentBFGBigExplosion:
		fallthrough
	case TentBossTeleport:
		fallthrough
	case TentPlainExplosion:
		fallthrough
	case TentChainFistSmoke:
		fallthrough
	case TentTrackerExplosion:
		fallthrough
	case TentTeleportEffect:
		fallthrough
	case TentDBallGoal:
		fallthrough
	case TentWidowSplash:
		fallthrough
	case TentNukeBlast:
		te.Position1 = m.ReadPosition()
	case TentParasiteAttack:
		fallthrough
	case TentMedicCableAttack:
		fallthrough
	case TentHeatBeam:
		fallthrough
	case TentMonsterHeatBeam:
		te.Entity1 = int16(m.ReadShort())
		te.Position1 = m.ReadPosition()
		te.Position2 = m.ReadPosition()
		te.Offset = m.ReadPosition()
	case TentGrappleCable:
		te.Entity1 = int16(m.ReadShort())
		te.Position1 = m.ReadPosition()
		te.Position2 = m.ReadPosition()
		te.Offset = m.ReadPosition()
	case TentLightning:
		te.Entity1 = int16(m.ReadShort())
		te.Entity2 = int16(m.ReadShort())
		te.Position1 = m.ReadPosition()
		te.Position2 = m.ReadPosition()
	case TentFlashlight:
		te.Position1 = m.ReadPosition()
		te.Entity1 = int16(m.ReadShort())
	case TentForceWall:
		te.Position1 = m.ReadPosition()
		te.Position2 = m.ReadPosition()
		te.Color = m.ReadByte()
	case TentSteam:
		te.Entity1 = int16(m.ReadShort())
		te.Count = m.ReadByte()
		te.Position1 = m.ReadPosition()
		te.Direction = m.ReadDirection()
		te.Color = m.ReadByte()
		te.Entity2 = int16(m.ReadShort())
		if te.Entity1 != -1 {
			te.Time = m.ReadLong()
		}
	case TentWidowBeamOut:
		te.Entity1 = int16(m.ReadShort())
		te.Position1 = m.ReadPosition()
	default:
		fmt.Printf("bad temp entity: %d\n", te.Type)
	}

	if *cli_args.Verbose {
		fmt.Printf(" * Temporary Entity [%d]\n", te.Type)
	}

	return te
}

func (m *MessageBuffer) ParseMuzzleFlash() MuzzleFlash {
	mf := MuzzleFlash{}
	mf.Entity = m.ReadShort()
	mf.Weapon = m.ReadByte()

	if *cli_args.Verbose {
		fmt.Printf(" * Muzzle Flash [%d]\n", mf.Weapon)
	}
	return mf
}

func (m *MessageBuffer) ParseLayout() Layout {
	layout := Layout{}
	layout.Data = m.ReadString()

	if *cli_args.Verbose {
		fmt.Printf(" * Layout\n")
	}

	return layout
}

// 2 bytes for every item
func (m *MessageBuffer) ParseInventory() {
	// we don't actually care about this, just parsing it
	inv := [MaxItems]uint16{}
	for i := 0; i < MaxItems; i++ {
		inv[i] = m.ReadShort()
	}

	if *cli_args.Verbose {
		fmt.Printf(" * Inventory\n")
	}
}

func (m *MessageBuffer) ParseCenterPrint() CenterPrint {
	c := CenterPrint{}
	c.Data = m.ReadString()

	if *cli_args.Verbose {
		fmt.Printf(" * CenterPrint\n")
	}

	return c
}
