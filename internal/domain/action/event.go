package action

type (
	Event struct {
		Id         Id         `json:"id"`
		Owner      ProfileId  `json:"owner"`
		Location   LocationId `json:"location"`
		GoldReward int64      `json:"goldReward"`
		ExpReward  int64      `json:"expReward"`
	}
)
