package bid_entity

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	Id        string
	UserId    string
	AuctionId string
	Amount    float64
	Timestamp time.Time
}

// Cria um novo lance (Bid)
func CreateBid(userId, auctionId string, amount float64) (*Bid, *internal_error.InternalError) {
	bid := &Bid{
		Id:        uuid.New().String(),
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validate(); err != nil {
		return nil, err
	}

	return bid, nil
}

// Valida os parâmetros do lance
func (b *Bid) Validate() *internal_error.InternalError {
	if err := uuid.Validate(b.UserId); err != nil {
		return internal_error.NewBadRequestError("UserId is not a valid UUID")
	}
	if err := uuid.Validate(b.AuctionId); err != nil {
		return internal_error.NewBadRequestError("AuctionId is not a valid UUID")
	}
	if b.Amount <= 0 {
		return internal_error.NewBadRequestError("Amount must be greater than zero")
	}

	return nil
}

// Interface do Repositório de Lances
type BidEntityRepository interface {
	// Cria múltiplos lances
	CreateBid(
		ctx context.Context,
		bidEntities []Bid) *internal_error.InternalError

	// Busca lances por ID do leilão
	FindBidByAuctionId(
		ctx context.Context, auctionId string) ([]Bid, *internal_error.InternalError)

	// Busca o lance vencedor por ID do leilão
	FindWinningBidByAuctionId(
		ctx context.Context, auctionId string) (*Bid, *internal_error.InternalError)
}
