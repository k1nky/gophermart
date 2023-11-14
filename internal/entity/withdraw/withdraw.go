package withdraw

import (
	"time"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

type ID uint64

//go:generate easyjson withdraw.go
//easyjson:json
type Withdraw struct {
	ID          ID                `json:"-"`
	Sum         float32           `json:"sum"`
	Number      order.OrderNumber `json:"order"`
	ProcessedAt time.Time         `json:"processed_at"`
	UserID      user.ID           `json:"-"`
}
