package profile

import (
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xuuid"
	"net/http"
)

const (
	ContextKey xctx.ContextKey = "profile_id"
)

type (
	Id string

	OwnerId string

	LocationId string

	Profile struct {
		Id       Id
		Name     string
		Owner    OwnerId
		Level    int64
		Gold     int64
		Location LocationId
	}
)

func New(name string, owner string, startLocation LocationId) (*Profile, error) {
	if name == "" || owner == "" {
		return nil, xhttp.NewError("name and owner cannot be null or empty", http.StatusBadRequest)
	}

	return &Profile{
		Id:       Id(xuuid.Base32UUID()),
		Name:     name,
		Owner:    OwnerId(owner),
		Level:    1,
		Gold:     0,
		Location: startLocation,
	}, nil
}
