// List implements a goroutine-safe player list.
package list

import (
	"sync"

	"github.com/toqueteos/minero/server/player"
)

// List is a simple goroutine-safe player list.
type List struct {
	sync.RWMutex

	list map[string]*player.Player
}

func New() *List {
	return &List{
		list: make(map[string]*player.Player),
	}
}

// Len returns the number of online players
func (l *List) Len() int {
	l.RLock()
	defer l.RUnlock()
	return len(list)
}

// Get gets a player from the list.
func (l *List) Get(name string) *player.Player {
	l.RLock()
	defer l.RUnlock()
	return l.list[name]
}

// Add adds a player to the list.
func (l *List) Add(p *player.Player) {
	l.Lock()
	l.list[p.Name] = p
	l.Unlock()
}

// Rem removes a player from the list.
func (l *List) Rem(name string) {
	l.Lock()
	delete(l.list, name)
	l.Unlock()
}
