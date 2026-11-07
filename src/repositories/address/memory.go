package address

import (
	"log"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Implementation

type memoryAddressRepo struct {
	shared.MemoryRepo[domain.Address]
}

//* Constructor

func NewMemoryAddressRepository(im ...domain.Address) domain.AddressRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Address]()

	for _, m := range im {
		if err := store.Set(m.ID, m); err != nil {
			log.Fatal(err)
		}
	}

	return &memoryAddressRepo{
		shared.MemoryRepo[domain.Address]{
			Store: store,
		},
	}
}

// NewJSONAddressRepository returns an in-memory repository backed by a JSON file:
// data is loaded from path on startup and written back on every change.
func NewJSONAddressRepository(path string) domain.AddressRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Address](path)

	return &memoryAddressRepo{
		shared.MemoryRepo[domain.Address]{
			Store: store,
		},
	}
}
