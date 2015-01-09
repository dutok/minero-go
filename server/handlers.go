package server

import (
	"github.com/minero/minero/server/player"
)

const (
	ReasonPiratedGame     = "Failed to login: User not premium."
	ReasonPiratedGameLong = "Failed to login: User not premium.\n\nYou have to buy the game to play on this server.\n\nGo to http://www.minecraft.net/store"
)

type HandlerFunc func(*Server, *player.Player)

// Packet Id is only read for "Client -> Server" packets
var HandlerFor = map[byte]HandlerFunc{
	0x00: HandlerFunc(Handle00), // KeepAlive
	0x02: HandlerFunc(Handle02), // Handshake
	0x03: HandlerFunc(Handle03), // ChatMessage
	0x07: HandlerFunc(Handle07), // EntityInteract
	0x0A: HandlerFunc(Handle0A), // Player
	0x0B: HandlerFunc(Handle0B), // PlayerPos
	0x0C: HandlerFunc(Handle0C), // PlayerLook
	0x0D: HandlerFunc(Handle0D), // PlayerPosLook
	0x0E: HandlerFunc(Handle0E), // PlayerAction
	0x0F: HandlerFunc(Handle0F), // PlayerBlockPlace
	0x10: HandlerFunc(Handle10), // ItemHeldChange
	0x12: HandlerFunc(Handle12), // Animation
	0x13: HandlerFunc(Handle13), // EntityAction
	0x65: HandlerFunc(Handle65), // WindowClose
	0x66: HandlerFunc(Handle66), // WindowClick
	0x6A: HandlerFunc(Handle6A), // ConfirmTransaction
	0x6B: HandlerFunc(Handle6B), // CreativeInventoryAction
	0x6C: HandlerFunc(Handle6C), // EnchantItem
	0x82: HandlerFunc(Handle82), // SignUpdate
	0xCA: HandlerFunc(HandleCA), // PlayerAbilities
	0xCB: HandlerFunc(HandleCB), // TabComplete
	0xCC: HandlerFunc(HandleCC), // ClientSettings
	0xCD: HandlerFunc(HandleCD), // ClientStatuses
	0xFA: HandlerFunc(HandleFA), // PluginMessage
	0xFC: HandlerFunc(HandleFC), // EncryptionKeyResponse
	0xFE: HandlerFunc(HandleFE), // ServerListPing
	0xFF: HandlerFunc(HandleFF), // DisconnectKick
}
