package product

import (
	"log"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Implementation

type memoryProductRepo struct {
	shared.MemoryRepo[domain.Product]
}

//* Constructor

func NewMemoryProductRepository(im ...domain.Product) domain.ProductRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Product]()

	for _, m := range im {
		if err := store.Set(m.ID, m); err != nil {
			log.Fatal(err)
		}
	}

	return &memoryProductRepo{
		shared.MemoryRepo[domain.Product]{
			Store: store,
		},
	}
}

// NewJSONProductRepository returns an in-memory repository backed by a JSON file:
// data is loaded from path on startup and written back on every change.
func NewJSONProductRepository(path string) domain.ProductRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Product](path)

	return &memoryProductRepo{
		shared.MemoryRepo[domain.Product]{
			Store: store,
		},
	}
}
