package ievent

import (
	"context"
	"sync"
)

type (
	Orchestrator[T any] struct {
		channels []chan T
		mutex    *sync.RWMutex
	}
)

func NewOrchestrator[T any]() *Orchestrator[T] {
	return &Orchestrator[T]{
		channels: make([]chan T, 0),
		mutex:    &sync.RWMutex{},
	}
}

func (orchestrator *Orchestrator[T]) PublishEvent(ctx context.Context, event T) error {
	orchestrator.mutex.RLock()
	defer orchestrator.mutex.RUnlock()

	for _, cha := range orchestrator.channels {
		cha <- event
	}

	return nil
}

func (orchestrator *Orchestrator[T]) RegisterListener(ctx context.Context) chan T {
	orchestrator.mutex.Lock()
	defer orchestrator.mutex.Unlock()

	channel := make(chan T)
	orchestrator.channels = append(orchestrator.channels, channel)

	return channel
}

func (orchestrator *Orchestrator[T]) UnregisterListener(ctx context.Context, channel chan T) {
	orchestrator.mutex.Lock()
	defer orchestrator.mutex.Unlock()

	close(channel)

	for index, cha := range orchestrator.channels {
		if cha == channel {
			orchestrator.channels = append(orchestrator.channels[:index], orchestrator.channels[index+1:]...)
			return
		}
	}
}
