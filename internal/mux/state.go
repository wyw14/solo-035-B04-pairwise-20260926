package mux

import (
	"sort"
	"sync"

	"github.com/wyw14/solo-035/internal/model"
)

type State struct {
	mu       sync.RWMutex
	programs map[model.EventID]Program
	active   map[model.EventID]bool
}

func NewState() *State {
	return &State{programs: make(map[model.EventID]Program), active: make(map[model.EventID]bool)}
}

func (state *State) Put(program Program) error {
	if err := program.Validate(); err != nil {
		return err
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	state.programs[program.EventID] = program
	return nil
}

func (state *State) Activate(eventID model.EventID, revision model.RevisionID) bool {
	state.mu.Lock()
	defer state.mu.Unlock()
	program, ok := state.programs[eventID]
	if !ok || program.Revision != revision {
		return false
	}
	state.active[eventID] = true
	return true
}

func (state *State) Deactivate(eventID model.EventID) {
	state.mu.Lock()
	defer state.mu.Unlock()
	delete(state.active, eventID)
}

func (state *State) Get(eventID model.EventID) (Program, bool) {
	state.mu.RLock()
	defer state.mu.RUnlock()
	program, ok := state.programs[eventID]
	return program, ok
}

func (state *State) List() []Program {
	state.mu.RLock()
	defer state.mu.RUnlock()
	items := make([]Program, 0, len(state.programs))
	for _, program := range state.programs {
		items = append(items, program)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return items
}
