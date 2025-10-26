package action

const (
	ActionTopic = "action"
)

type (
	Event struct {
		Id         Id         `json:"id"`
		ProfileId  ProfileId  `json:"profileId"`
		Location   LocationId `json:"location"`
		GoldReward int64      `json:"goldReward"`
		ExpReward  int64      `json:"expReward"`
	}
)
