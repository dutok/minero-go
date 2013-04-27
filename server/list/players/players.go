// Package players implements a goroutine-safe player list.
package players

import (
	"sync"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// PlayersList is a simple goroutine-safe player list.
type PlayersList struct {
	sync.RWMutex

	list map[string]*player.Player
}

func New() PlayersList {
	return PlayersList{
		list: make(map[string]*player.Player),
	}
}

// Len returns the number of online players
func (l PlayersList) Len() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.list)
}

// GetPlayer gets a player from the list by his/her name.
func (l PlayersList) GetPlayer(name string) *player.Player {
	l.RLock()
	defer l.RUnlock()
	return l.list[name]
}

// AddPlayer adds a player to the list.
func (l PlayersList) AddPlayer(p *player.Player) {
	l.Lock()
	l.list[p.Name] = p
	l.Unlock()
}

// RemPlayer removes a player from the list.
func (l PlayersList) RemPlayer(p *player.Player) {
	l.Lock()
	delete(l.list, p.Name)
	l.Unlock()
}

// BroadcastPacket sends a packet to all online players.
func (l PlayersList) BroadcastPacket(pkt packet.Packet) {
	l.RLock()
	for _, p := range l.list {
		if p.Ready {
			pkt.WriteTo(p.Conn)
		}
	}
	l.RUnlock()
}

// BroadcastMessage send a message to all online players.
func (l PlayersList) BroadcastMessage(msg string) {
	l.RLock()
	for _, p := range l.list {
		if p.Ready {
			p.SendMessage(msg)
		}
	}
	l.RUnlock()
}

// BroadcastLogin initializes all previously online clients to a new player.
func (l PlayersList) BroadcastLogin(to *player.Player) {
	l.RLock()
	for _, p := range l.list {
		if p.Ready {
			r := &packet.EntityNamedSpawn{
				Entity:   p.Id(),
				Name:     p.Name,
				X:        p.X,
				Y:        p.Y,
				Z:        p.Z,
				Yaw:      p.Yaw,
				Pitch:    p.Pitch,
				Item:     0,
				Metadata: player.JustLoginMetadata(p.Name),
			}
			r.WriteTo(to.Conn)
		}
	}
	l.RUnlock()
}
