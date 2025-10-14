package action

import (
	"github.com/VaynerAkaWalo/go-toolkit/xuuid"
	"math/rand"
)

type (
	Id string

	Action struct {
		Id       Id
		Player   ProfileId
		Location LocationId
		Reward   int64
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

func calculateReward(base int64, bonusMultiplier float64) int64 {
	sample := rand.Float64()

	if sample <= bonusMultiplier {
		return base + 1
	}

	return base
}
