package minecraft

const (
	PacketKeepAlive          byte = iota
	PacketLoginInfo               // 0x01
	PacketHandshake               // 0x02
	PacketChatMessage             // 0x03
	PacketTimeUpdate              // 0x04
	PacketEntityEquipment         // 0x05
	PacketSpawnPosition           // 0x06
	PacketEntityInteract          // 0x07
	PacketHealthUpdate            // 0x08
	PacketRespawn                 // 0x09
	PacketPlayer                  // 0x0a
	PacketPlayerPos               // 0x0b
	PacketPlayerLook              // 0x0c
	PacketPlayerPosLook           // 0x0d
	PacketPlayerAction            // 0x0e
	PacketPlayerBlockPlace        // 0x0f
	PacketItemHeldChange          // 0x10
	PacketBedUse                  // 0x11
	PacketAnimation               // 0x12
	PacketEntityAction            // 0x13
	PacketEntityNamedSpawn        // 0x14
	PacketItemCollect             // 0x15
	PacketSpawnObjectVehicle      // 0x16
	PacketSpawnMob                // 0x17
	PacketSpawnPainting           // 0x18
	PacketSpawnExperienceOrb      // 0x19
	PacketEntityVelocity          // 0x1a
	PacketEntityDestroy           // 0x1b
	PacketEntity                  // 0x1c
	PacketEntityRelMove           // 0x1d
	PacketEntityLook              // 0x1e
	PacketEntityLookRelMove       // 0x1f
	PacketEntityTeleport          // 0x20
	PacketEntityHeadLook          // 0x21
	PacketEntityStatus            // 0x22
	PacketEntityAttach            // 0x23
	PacketEntityMetadata          // 0x24
	PacketEntityEffect            // 0x25
	PacketEntityEffectRemove      // 0x26
	PacketSetExperience           // 0x27

	PacketChunkData           byte = iota + 0x33
	PacketBlockChangeMulti         // 0x34
	PacketBlockChange              // 0x35
	PacketBlockAction              // 0x36
	PacketBlockBreakAnimation      // 0x37
	PacketMapChunkBulk             // 0x38
	PacketExplosion                // 0x39
	PacketSoundEffect              // 0x40
	PacketSoundEffectNamed         // 0x41
	PacketParticle                 // 0x42

	PacketGameStateChange   byte = iota + 0x46
	PacketEntityGlobalSpawn      // 0x47

	PacketWindowOpen              byte = iota + 0x64
	PacketWindowClose                  // 0x65
	PacketWindowClick                  // 0x66
	PacketWindowSlotSet                // 0x67
	PacketWindowSetItems               // 0x68
	PacketWindowUpdateProperty         // 0x69
	PacketConfirmTransaction           // 0x6a
	PacketCreativeInventoryAction      // 0x6b
	PacketEnchantItem                  // 0x6c

	PacketSignUpdate       byte = iota + 0x82
	PacketItemData              // 0x83
	PacketTileEntityUpdate      // 0x84

	PacketStatIncrement     byte = iota + 0xC8
	PacketPlayerTabListPing      // 0xc9
	PacketPlayerAbilities        // 0xca
	PacketTabComplete            // 0xcb
	PacketClientSettings         // 0xcc
	PacketClientStatuses         // 0xcd
	PacketScoreObjective         // 0xce
	PacketScoreUpdate            // 0xcf
	PacketScoreDisplay           // 0xd0
	PacketScoreTeams             // 0xd1

	PacketPluginMessage byte = 0xFA

	PacketEncryptionKeyResponse byte = iota + 0xFC
	PacketEncryptionKeyRequest       // 0xFD
	PacketServerListPing             // 0xFE
	PacketDisconnect                 // 0xFF
)

var packetNames = map[byte]string{
	0x00: "KeepAlive",
	0x01: "LoginInfo",
	0x02: "Handshake",
	0x03: "ChatMessage",
	0x04: "TimeUpdate",
	0x05: "EntityEquipment",
	0x06: "SpawnPosition",
	0x07: "EntityInteract",
	0x08: "HealthUpdate",
	0x09: "Respawn",
	0x0A: "Player",
	0x0B: "PlayerPos",
	0x0C: "PlayerLook",
	0x0D: "PlayerPosLook",
	0x0E: "PlayerAction",
	0x0F: "PlayerBlockPlace",
	0x10: "ItemHeldChange",
	0x11: "BedUse",
	0x12: "Animation",
	0x13: "EntityAction",
	0x14: "EntityNamedSpawn",
	0x16: "ItemCollect",
	0x17: "SpawnObjectVehicle",
	0x18: "SpawnMob",
	0x19: "SpawnPainting",
	0x1A: "SpawnExperienceOrb",
	0x1C: "EntityVelocity",
	0x1D: "EntityDestroy",
	0x1E: "Entity",
	0x1F: "EntityRelMove",
	0x20: "EntityLook",
	0x21: "EntityLookRelMove",
	0x22: "EntityTeleport",
	0x23: "EntityHeadLook",
	0x26: "EntityStatus",
	0x27: "EntityAttach",
	0x28: "EntityMetadata",
	0x29: "EntityEffect",
	0x2A: "EntityEffectRemove",
	0x2B: "SetExperience",
	0x33: "ChunkData",
	0x34: "BlockChangeMulti",
	0x35: "BlockChange",
	0x36: "BlockAction",
	0x37: "BlockBreakAnimation",
	0x38: "MapChunkBulk",
	0x3C: "Explosion",
	0x3D: "SoundEffect",
	0x3E: "SoundEffectNamed",
	0x3F: "Particle",
	0x46: "GameStateChange",
	0x47: "EntityGlobalSpawn",
	0x64: "WindowOpen",
	0x65: "WindowClose",
	0x66: "WindowClick",
	0x67: "WindowSlotSet",
	0x68: "WindowSetItems",
	0x69: "WindowUpdateProperty",
	0x6A: "ConfirmTransaction",
	0x6B: "CreativeInventoryAction",
	0x6C: "EnchantItem",
	0x82: "SignUpdate",
	0x83: "ItemData",
	0x84: "TileEntityUpdate",
	0xC8: "StatIncrement",
	0xC9: "PlayerTabListPing",
	0xCA: "PlayerAbilities",
	0xCB: "TabComplete",
	0xCC: "ClientSettings",
	0xCD: "ClientStatuses",
	0xCE: "ScoreboardObjective",
	0xCF: "ScoreUpdate",
	0xD0: "ScoreboardDisplay",
	0xD1: "Teams",
	0xFA: "PluginMessage",
	0xFC: "EncryptionKeyResponse",
	0xFD: "EncryptionKeyRequest",
	0xFE: "ServerListPing",
	0xFF: "DisconnectKick",
}
