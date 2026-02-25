package events

import (
	"server/internal/core/shared/types"
	"sync"
)

type events struct {
	mu           sync.RWMutex
	evtListeners types.Event
}

func New() types.IEvents {
	return &events{
		evtListeners: types.Event{},
	}
}

func (e *events) On(name string, listener ...types.EventListener) {
	e.addListener(name, listener...)
}

func (e *events) Emit(name string, data ...any) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.evtListeners != nil {
		if listeners := e.evtListeners[name]; len(listeners) > 0 {
			for key := range listeners {
				lis := listeners[key]
				if lis != nil {
					lis(data...)
				}
			}
		}
	}
}

func (e *events) addListener(name string, listener ...types.EventListener) {
	if len(listener) == 0 {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.evtListeners == nil {
		e.evtListeners = types.Event{}
	}

	listeners := e.evtListeners[name]

	if listeners == nil {
		listeners = []types.EventListener{}
	}

	e.evtListeners[name] = append(listeners, listener...)
}
