package action

import (
	"github.com/VaynerAkaWalo/go-toolkit/xuuid"
	"math/rand"
)

type (
	Id string

	Action struct {
		Id       Id         `json:"id"`
		Player   ProfileId  `json:"player"`
		Location LocationId `json:"location"`
		Reward   int64      `json:"reward"`
	}
)

func New(profile Profile, location Location) *Action {
	return &Action{
		Id:       Id(xuuid.Base32UUID()),
		Player:   profile.Id,
		Location: location.Id,
		Reward:   calculateReward(1, location.Multiplier),
	}
}

func (action Action) CreateEvent() Event {
	return Event{
		Id:         action.Id,
		Owner:      action.Player,
		Location:   action.Location,
		GoldReward: action.Reward,
		ExpReward:  1,
	}
}

func calculateReward(base int64, bonusMultiplier float64) int64 {
	sample := rand.Float64()

	if sample <= bonusMultiplier {
		return base + 1
	}

	return base
}
