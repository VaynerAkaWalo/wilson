package aprofile

import (
	"context"
	"golang-template/internal/domain/profile"
	"sync"
)

type (
	InMemoryProfileStore struct {
		profiles map[profile.OwnerId][]profile.Profile
		mutex    *sync.RWMutex
	}
)

func NewRepository() profile.Repository {
	return InMemoryProfileStore{
		profiles: make(map[profile.OwnerId][]profile.Profile),
		mutex:    &sync.RWMutex{},
	}
}

func (store InMemoryProfileStore) GetProfilesByOwner(ctx context.Context, id profile.OwnerId) ([]profile.Profile, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	return store.profiles[id], nil
}

func (store InMemoryProfileStore) Save(ctx context.Context, newProfile *profile.Profile) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	array, found := store.profiles[newProfile.Owner]
	if !found {
		array = []profile.Profile{*newProfile}
	} else {
		array = append(array, *newProfile)
	}

	store.profiles[newProfile.Owner] = array
	return nil
}
