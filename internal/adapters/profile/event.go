package adapter_profile

type (
	EventType string

	ActionEvent struct {
		Id   string `json:"id"`
		Gold int64  `json:"gold"`
		Exp  int64  `json:"exp"`
		Date int64  `json:"date"`
	}

	GoldChangeEvent struct {
		Id   string `json:"id"`
		Gold int64  `json:"gold"`
	}
)

const (
	NoEvent    EventType = "no-event"
	Action     EventType = "action"
	GoldChange EventType = "gold-change"
)
