package action

import "fmt"

type (
	ProfileId string

	Profile struct {
		Id       ProfileId
		Location LocationId
		Level    int64
		Gold     int64
	}
)

func (profile *Profile) ConsumeAction(action Action) error {
	if action.Player != profile.Id {
		return fmt.Errorf("cannot consume action owned by another player")
	}

	profile.Gold += action.Reward
	return nil
}
