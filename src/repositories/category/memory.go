package category

import (
	"log"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Implementation

type memoryCategoryRepo struct {
	shared.MemoryRepo[domain.Category]
}

//* Constructor

func NewMemoryCategoryRepository(im ...domain.Category) domain.CategoryRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Category]()

	for _, m := range im {
		if err := store.Set(m.ID, m); err != nil {
			log.Fatal(err)
		}
	}

	return &memoryCategoryRepo{
		shared.MemoryRepo[domain.Category]{
			Store: store,
		},
	}
}

// NewJSONCategoryRepository returns an in-memory repository backed by a JSON file:
// data is loaded from path on startup and written back on every change.
func NewJSONCategoryRepository(path string) domain.CategoryRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Category](path)

	return &memoryCategoryRepo{
		shared.MemoryRepo[domain.Category]{
			Store: store,
		},
	}
}
