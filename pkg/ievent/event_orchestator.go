package ievent

import (
	"context"
	"sync"
)

type (
	Orchestrator struct {
		channels []chan interface{}
		mutex    *sync.RWMutex
	}
)

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		channels: []chan interface{}{},
		mutex:    &sync.RWMutex{},
	}
}

func (o *Orchestrator) PublishEvent(ctx context.Context, event interface{}) error {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	for _, cha := range o.channels {
		cha <- event
	}

	return nil
}

func (o *Orchestrator) RegisterListener(ctx context.Context) (chan interface{}, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	channel := make(chan interface{})
	o.channels = append(o.channels, channel)

	return channel, nil
}

func (o *Orchestrator) UnregisterListener(ctx context.Context, channel chan interface{}) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	close(channel)

	for index, cha := range o.channels {
		if cha == channel {
			o.channels = append(o.channels[:index], o.channels[index+1:]...)
			return
		}
	}
}
