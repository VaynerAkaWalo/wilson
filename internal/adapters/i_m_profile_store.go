package adapters

import (
	"context"
	"fmt"
	"github.com/VaynerAkaWalo/go-toolkit/xuuid"
	"golang-template/internal/domain/action"
	"golang-template/internal/domain/profile"
	"golang-template/internal/domain/transaction"
	"sync"
)

type (
	Profile struct {
		Id          string
		Name        string
		Owner       string
		Location    string
		Level       int64
		Gold        int64
		GoldVersion int64
	}

	InMemoryProfileStore struct {
		profiles []Profile
		mutex    *sync.RWMutex
	}
)

func NewRepository() *InMemoryProfileStore {
	return &InMemoryProfileStore{
		profiles: make([]Profile, 0),
		mutex:    &sync.RWMutex{},
	}
}

func (store *InMemoryProfileStore) GetProfilesByOwner(ctx context.Context, id profile.OwnerId) ([]profile.Profile, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, prof := range store.profiles {
		if prof.Owner == string(id) {
			return []profile.Profile{*prof.toProfileView()}, nil
		}
	}

	return []profile.Profile{}, nil
}

func (store *InMemoryProfileStore) Save(ctx context.Context, p *profile.Profile) error {
	newProfile := Profile{
		Id:       string(p.Id),
		Name:     p.Name,
		Owner:    string(p.Owner),
		Location: xuuid.Base32UUID(),
		Level:    p.Level,
		Gold:     p.Gold,
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.profiles = append(store.profiles, newProfile)

	return nil
}

func (store *InMemoryProfileStore) Get(ctx context.Context, id action.ProfileId) (action.Profile, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, prof := range store.profiles {
		if prof.Id == string(id) {
			return *prof.toActionView(), nil
		}
	}

	return action.Profile{}, fmt.Errorf("profile not found")
}

func (store *InMemoryProfileStore) GetAll(ctx context.Context) ([]action.Profile, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	profiles := make([]action.Profile, 0)

	for _, prof := range store.profiles {
		profiles = append(profiles, *prof.toActionView())
	}

	return profiles, nil
}

func (p *Profile) toProfileView() *profile.Profile {
	return &profile.Profile{
		Id:    profile.Id(p.Id),
		Name:  p.Name,
		Owner: profile.OwnerId(p.Owner),
		Level: p.Level,
		Gold:  p.Gold,
	}
}

func (p *Profile) toActionView() *action.Profile {
	return &action.Profile{
		Id:       action.ProfileId(p.Id),
		Location: action.LocationId(p.Location),
	}
}

func (p *Profile) toBalanceView() *transaction.Balance {
	return &transaction.Balance{
		Profile: p.Id,
		Gold:    p.Gold,
		Version: p.GoldVersion,
	}
}

func (store *InMemoryProfileStore) GetBalance(ctx context.Context, profile string) (transaction.Balance, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, prof := range store.profiles {
		if prof.Id == profile {
			return *prof.toBalanceView(), nil
		}
	}

	return transaction.Balance{}, fmt.Errorf("profile not found")
}

func (store *InMemoryProfileStore) UpdateBalance(ctx context.Context, balance transaction.Balance) (transaction.Balance, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for i := range store.profiles {
		prof := &store.profiles[i]

		if prof.Id == balance.Profile {
			if balance.Version != prof.GoldVersion {
				return transaction.Balance{}, transaction.VersionMismatchError{}
			}

			prof.Gold = balance.Gold
			prof.GoldVersion++

			return *prof.toBalanceView(), nil
		}
	}

	return transaction.Balance{}, fmt.Errorf("profile not found")
}
