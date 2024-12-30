package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDatabase(t *testing.T) *mongo.Database {
	// Implement the setup logic for the test database here
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client.Database("testdb")
}

func TestCloseExpiredAuctions(t *testing.T) {
	ctx := context.Background()
	database := setupTestDatabase(t)
	repo := NewAuctionRepository(database)

	// Crie um leilão expirado
	expiredAuction := &AuctionEntityMongo{
		Id:          "test-expired",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   1,
		Status:      0, // Active
		Timestamp:   time.Now().Add(-2 * time.Hour).Unix(),
	}
	repo.Collection.InsertOne(ctx, expiredAuction)

	// Executa o fechamento
	err := repo.CloseExpiredAuctions(ctx, 1*time.Second)
	assert.Nil(t, err)

	// Verifica se foi fechado
	var result AuctionEntityMongo
	repo.Collection.FindOne(ctx, bson.M{"_id": "test-expired"}).Decode(&result)
	assert.Equal(t, auction_entity.Closed, result.Status)
}
