package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func (ar *AuctionRepository) StartAuctionExpirationWatcher(ctx context.Context) {
	durationStr := os.Getenv("AUCTION_INTERVAL")
	auctionDuration, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Error("Invalid AUCTION_INTERVAL in .env", err)
		return
	}

	ticker := time.NewTicker(auctionDuration / 2) // Verificação periódica
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := ar.CloseExpiredAuctions(ctx, auctionDuration)
			if err != nil {
				logger.Error("Error while closing expired auctions", err)
			}
		case <-ctx.Done():
			logger.Info("Auction expiration watcher stopped.")
			return
		}
	}
}

func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context, auctionDuration time.Duration) *internal_error.InternalError {
	filter := bson.M{
		"status": auction_entity.Active,
		"timestamp": bson.M{
			"$lt": time.Now().Add(-auctionDuration).Unix(),
		},
	}
	update := bson.M{
		"$set": bson.M{
			"status": auction_entity.Closed,
		},
	}

	_, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to close expired auctions", err)
		return internal_error.NewInternalServerError("Failed to close expired auctions")
	}

	logger.Info("Expired auctions closed successfully.")
	return nil
}
