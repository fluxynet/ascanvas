package memory

import (
	"context"
	"sync"

	"github.com/fluxynet/ascanvas"
)

type Memory struct {
	listeners map[string]map[int]chan ascanvas.CanvasEvent
	mutex     sync.Mutex
	counter   int
}

func New() *Memory {
	var listeners = make(map[string]map[int]chan ascanvas.CanvasEvent)

	return &Memory{
		listeners: listeners,
	}
}

func (m *Memory) Observe(ctx context.Context, id string) (ascanvas.StopObserveFunc, <-chan ascanvas.CanvasEvent, error) {
	var c = make(chan ascanvas.CanvasEvent)

	defer m.mutex.Unlock()
	m.mutex.Lock()

	if _, ok := m.listeners[id]; !ok {
		m.listeners[id] = make(map[int]chan ascanvas.CanvasEvent)
	}

	m.counter += 1
	var i = m.counter
	m.listeners[id][i] = c

	var stop = func() {
		defer m.mutex.Unlock()
		m.mutex.Lock()
		if _, ok := m.listeners[id]; ok {
			delete(m.listeners[id], i)
		}
	}

	return stop, c, nil
}

func (m *Memory) Broadcast(ctx context.Context, event ascanvas.CanvasEvent) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	var n = len(m.listeners[event.Canvas.Id]) + len(m.listeners[ascanvas.ObserveALL])

	var wg sync.WaitGroup
	wg.Add(n)

	for i := range m.listeners[event.Canvas.Id] {
		go func(l chan ascanvas.CanvasEvent) {
			defer wg.Done()
			l <- event
		}(m.listeners[event.Canvas.Id][i])
	}

	for i := range m.listeners[ascanvas.ObserveALL] {
		go func(l chan ascanvas.CanvasEvent) {
			defer wg.Done()
			l <- event
		}(m.listeners[ascanvas.ObserveALL][i])
	}

	wg.Wait()

	return nil
}

func (m *Memory) Close() error {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	for i := range m.listeners {
		for j := range m.listeners[i] {
			close(m.listeners[i][j])
		}
	}

	m.listeners = nil

	return nil
}
