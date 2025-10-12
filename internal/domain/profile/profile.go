package profile

import (
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xuuid"
	"net/http"
)

type (
	Id string

	OwnerId string

	Profile struct {
		Id    Id
		Name  string
		Owner OwnerId
		Level int64
		Gold  int64
	}
)

func New(name string, owner string) (*Profile, error) {
	if name == "" || owner == "" {
		return nil, xhttp.NewError("name and owner cannot be null or empty", http.StatusBadRequest)
	}

	return &Profile{
		Id:    Id(xuuid.Base32UUID()),
		Name:  name,
		Owner: OwnerId(owner),
		Level: 1,
		Gold:  0,
	}, nil
}
