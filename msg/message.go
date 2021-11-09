package msg

import (
	//"bytes"
	//"encoding/binary"
	//"fmt"

    //"github.com/packetflinger/q2demo/msg"
)

type MessageBuffer struct {
	Buffer []byte
	Index  int
	Length int // maybe not needed
}

type ServerData struct {
	Protocol        int32
	ServerCount     int32
	Demo            int8
	GameDir         string
	ClientNumber    int16
	MapName         string
}

type ConfigString struct {
    Index int16
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
    PlayerMove  PlayerMoveState  
    ViewAngles  [3]int16
    ViewOffset  [3]int8
    KickAngles  [3]int8
    GunAngles   [3]int8
    GunOffset   [3]int8
    GunIndex    uint8
    GunFrame    uint8
    Blend       [4]uint8
    FOV         uint8
    RDFlags     uint8
    Stats       [32]int16
}

type Frame struct {
    Number      int32
    Delta       int32
    Suppressed  int8
    AreaBytes   int8
    AreaBits    []byte
}


func ParseServerData(m *MessageBuffer) ServerData {
    sd := ServerData{}

    sd.Protocol = ReadLong(m)
    sd.ServerCount = ReadLong(m)
    sd.Demo = int8(ReadByte(m))
    sd.GameDir = ReadString(m)
    sd.ClientNumber = int16(ReadShort(m))
    sd.MapName = ReadString(m)

    return sd
}

func ParseConfigString(m *MessageBuffer) ConfigString {
    cs := ConfigString{
        Index: int16(ReadShort(m)),
        String: ReadString(m),
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
    bits := uint32(ReadByte(m))
    
    if bits & EntityMoreBits1 != 0 {
        bits |= (uint32(ReadByte(m)) << 8)
    }

    if bits & EntityMoreBits2 != 0 {
        bits |= (uint32(ReadByte(m)) << 16)
    }

    if bits & EntityMoreBits3 != 0 {
        bits |= (uint32(ReadByte(m)) << 24)
    }

    return uint32(bits)
}

func ParseEntityNumber(m *MessageBuffer, flags uint32) uint16 {
    num := uint16(0)
    if flags & EntityNumber16 != 0 {
        num = uint16(ReadShort(m))
    } else {
        num = uint16(ReadByte(m))
    }

    return num
}

func ParseEntity(m *MessageBuffer, from PackedEntity, num uint16, bits uint32) PackedEntity {
    to := from
    to.Number = uint32(num)
    
    if bits == 0 {
        return to
    }

    if bits & EntityModel != 0 {
        to.ModelIndex = uint8(ReadByte(m))
    }

    if bits & EntityModel2 != 0 {
        to.ModelIndex2 = uint8(ReadByte(m))
    }

    if bits & EntityModel3 != 0 {
        to.ModelIndex3 = uint8(ReadByte(m))
    }

    if bits & EntityModel4 != 0 {
        to.ModelIndex4 = uint8(ReadByte(m))
    }

    if bits & EntityFrame8 != 0 {
        to.Frame = uint16(ReadByte(m))
    }

    if bits & EntityFrame16 != 0 {
        to.Frame = uint16(ReadShort(m))
    }

    if (bits & (EntitySkin8 | EntitySkin16)) == (EntitySkin8 | EntitySkin16) {
        to.SkinNum = uint32(ReadLong(m))
    } else if bits & EntitySkin8 != 0 {
        to.SkinNum = uint32(ReadByte(m))
    } else if bits & EntitySkin16 != 0 {
        to.SkinNum = uint32(ReadWord(m))
    }

    if (bits & (EntityEffects8 | EntityEffects16)) == (EntityEffects8 | EntityEffects16) {
        to.Effects = uint32(ReadLong(m))
    } else if bits & EntityEffects8 != 0 {
        to.Effects = uint32(ReadByte(m))
    } else if bits & EntityEffects16 != 0 {
        to.Effects = uint32(ReadWord(m))
    }

    if (bits & (EntityRenderFX8 | EntityRenderFX16)) == (EntityRenderFX8 | EntityRenderFX16) {
        to.RenderFX = uint32(ReadLong(m))
    } else if bits & EntityRenderFX8 != 0 {
        to.RenderFX = uint32(ReadByte(m))
    } else if bits & EntityRenderFX16 != 0 {
        to.RenderFX = uint32(ReadWord(m))
    }

    if bits & EntityOrigin1 != 0 {
        to.Origin[0] = int16(ReadShort(m))
    }

    if bits & EntityOrigin2 != 0 {
        to.Origin[1] = int16(ReadShort(m))
    }

    if bits & EntityOrigin3 != 0 {
        to.Origin[2] = int16(ReadShort(m))
    }

    if bits & EntityAngle1 != 0 {
        to.Angles[0] = int16(ReadByte(m))
    }

    if bits & EntityAngle2 != 0 {
        to.Angles[1] = int16(ReadByte(m))
    }

    if bits & EntityAngle3 != 0 {
        to.Angles[2] = int16(ReadByte(m))
    }

    if bits & EntityOldOrigin != 0 {
        to.OldOrigin[0] = int16(ReadShort(m))
        to.OldOrigin[1] = int16(ReadShort(m))
        to.OldOrigin[2] = int16(ReadShort(m))
    }

    if bits & EntitySound != 0 {
        to.Sound = uint8(ReadByte(m))
    }

    if bits & EntityEvent != 0 {
        to.Event = uint8(ReadByte(m))
    }

    if bits & EntitySolid != 0 {
        to.Solid = uint32(ReadWord(m))
    }

    return to
}

func ParseStuffText(m *MessageBuffer) StuffText {
    str := StuffText{String: ReadString(m)}
    return str
}

func ParseFrame(m *MessageBuffer) Frame {
    N := int32(ReadLong(m))
    D := int32(ReadLong(m))
    S := int8(ReadByte(m))
    A := int8(ReadByte(m))
    Ab := ReadData(m, int(A))
    
    fr := Frame{
        Number: N,
        Delta: D,
        Suppressed: S,
        AreaBytes: A,
        AreaBits: Ab,
    }

    return fr
}

func ParsePlayerstate(m *MessageBuffer) PackedPlayer {
    bits := ReadWord(m)
    pm := PlayerMoveState{}
    ps := PackedPlayer{}
    
    if bits & PlayerType != 0 {
        pm.Type = int8(ReadByte(m))
    }

    if bits & PlayerOrigin != 0 {
        pm.Origin[0] = int16(ReadShort(m))
        pm.Origin[1] = int16(ReadShort(m))
        pm.Origin[2] = int16(ReadShort(m))
    }

    if bits & PlayerVelocity != 0 {
        pm.Velocity[0] = int16(ReadShort(m))
        pm.Velocity[1] = int16(ReadShort(m))
        pm.Velocity[2] = int16(ReadShort(m))
    }

    if bits & PlayerTime != 0 {
        pm.Time = ReadByte(m)
    }

    if bits & PlayerFlags != 0 {
        pm.Flags = ReadByte(m)
    }

    if bits & PlayerGravity != 0 {
        pm.Gravity = int16(ReadShort(m))
    }

    if bits & PlayerDeltaAngles != 0 {
        pm.DeltaAngles[0] = int16(ReadShort(m))
        pm.DeltaAngles[1] = int16(ReadShort(m))
        pm.DeltaAngles[2] = int16(ReadShort(m))
    }

    if bits & PlayerViewOffset != 0 {
        ps.ViewOffset[0] = int8(ReadChar(m))
        ps.ViewOffset[1] = int8(ReadChar(m))
        ps.ViewOffset[2] = int8(ReadChar(m))
    }

    if bits & PlayerViewAngles != 0 {
        pm.ViewAngles[0] = int16(ReadShort(m))
        pm.ViewAngles[1] = int16(ReadShort(m))
        pm.ViewAngles[2] = int16(ReadShort(m))
    }

    if bits & PlayerWeaponIndex != 0 {
        ps.GunIndex = int8(ReadByte(m))
    }

    if bits & PlayerWeaponFrame != 0 {
        ps.GunFrame = int8(ReadByte(m))
        ps.GunOffset[0] = int8(ReadChar(m))
        ps.GunOffset[1] = int8(ReadChar(m))
        ps.GunOffset[2] = int8(ReadChar(m))
        ps.GunAngles[0] = int8(ReadChar(m))
        ps.GunAngles[1] = int8(ReadChar(m))
        ps.GunAngles[2] = int8(ReadChar(m))
    }

    if bits & PlayerBlend != 0 {
        ps.Blend[0] = int8(ReadChar(m))
        ps.Blend[1] = int8(ReadChar(m))
        ps.Blend[2] = int8(ReadChar(m))
        ps.Blend[3] = int8(ReadChar(m))
    }

    if bits & PlayerFOV != 0 {
        ps.FOV = int8(ReadByte(m))
    }

    if bits & PlayerRDFlags != 0 {
        ps.RDFlags = int8(ReadByte(m))
    }

    statbits := int32(ReadLong(m))
    for i:=0; i<32; i++ {
        if statbits & (1<<i) != 0 {
            ps.Stats[i] = uint16(ReadShort(m))
        }
    }

    ps.PlayerMoveState = pm
}