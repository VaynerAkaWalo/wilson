package transaction

import (
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xuuid"
	"net/http"
	"time"
)

const (
	ContextKey xctx.ContextKey = "transaction_id"
)

type (
	Transaction struct {
		Id            Id     `json:"id"`
		Profile       string `json:"profile"`
		BalanceChange int64  `json:"balanceChange"`
		CreatedAt     int64  `json:"createdAt"`
	}

	Id string
)

func New(profile string, balanceChange int64) (Transaction, error) {
	if balanceChange <= 0 {
		return Transaction{}, xhttp.NewError("balance change cannot be negative", http.StatusBadRequest)
	}

	return Transaction{
		Id:            Id(xuuid.UUID()),
		Profile:       profile,
		BalanceChange: balanceChange,
		CreatedAt:     time.Now().Unix(),
	}, nil
}
