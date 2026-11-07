package order

import (
	"log"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Implementation

type memoryOrderRepo struct {
	shared.MemoryRepo[domain.Order]
}

//* Constructor

func NewMemoryOrderRepository(im ...domain.Order) domain.OrderRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Order]()

	for _, m := range im {
		if err := store.Set(m.ID, m); err != nil {
			log.Fatal(err)
		}
	}

	return &memoryOrderRepo{
		shared.MemoryRepo[domain.Order]{
			Store: store,
		},
	}
}

// NewJSONOrderRepository returns an in-memory repository backed by a JSON file:
// data is loaded from path on startup and written back on every change.
func NewJSONOrderRepository(path string) domain.OrderRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Order](path)

	return &memoryOrderRepo{
		shared.MemoryRepo[domain.Order]{
			Store: store,
		},
	}
}
