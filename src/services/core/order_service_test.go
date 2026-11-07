package core

import (
	"testing"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/address"
	"github.com/ZaphCode/clean-arch/src/repositories/order"
	"github.com/ZaphCode/clean-arch/src/repositories/product"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type OrderServiceSuite struct {
	suite.Suite
	service *orderService
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceSuite))
}

func (s *OrderServiceSuite) SetupSuite() {
	s.T().Logf("\n-------------- init ---------------")

	ordRepo := order.NewMemoryOrderRepository()
	addrRepo := address.NewMemoryAddressRepository()

	s.service = &orderService{
		ordRepo:  ordRepo,
		addrRepo: addrRepo,
	}
}

//* Tests

func (s *OrderServiceSuite) TestCreateOrder() {
	ord := domain.Order{}
	s.service.Create(&ord)
}

func (s *OrderServiceSuite) TestGetDetailedByID() {
	ordRepo := order.NewMemoryOrderRepository()
	svc := &orderService{
		ordRepo:  ordRepo,
		addrRepo: address.NewMemoryAddressRepository(utils.AddrExp1),
		prodRepo: product.NewMemoryProductRepository(utils.ProductExp1),
	}

	owner := uuid.New()

	// Points to an existing address and product.
	complete := domain.Order{Model: domain.Model{ID: uuid.New()}, UserID: owner, AddressID: utils.AddrExp1.ID,
		Products: []domain.OrderProduct{{ID: utils.ProductExp1.ID, Quantity: 1}}}
	// Points to an address and a product that no longer exist.
	orphan := domain.Order{Model: domain.Model{ID: uuid.New()}, UserID: owner, AddressID: uuid.New(),
		Products: []domain.OrderProduct{{ID: utils.ProductExp1.ID, Quantity: 1}, {ID: uuid.New(), Quantity: 2}}}

	s.Require().NoError(ordRepo.Save(&complete))
	s.Require().NoError(ordRepo.Save(&orphan))

	s.Run("complete order", func() {
		got, err := svc.GetDetailedByID(complete.ID, owner)
		s.Require().NoError(err)
		s.Require().NotNil(got)
		s.Equal(utils.AddrExp1.ID, got.Address.ID)
		s.Require().Len(got.Products, 1)
		s.Equal(utils.ProductExp1.ID, got.Products[0].ID)
		s.Equal(uint(1), got.Products[0].Quantity)
	})

	s.Run("missing address and product don't panic", func() {
		got, err := svc.GetDetailedByID(orphan.ID, owner)
		s.Require().NoError(err)
		s.Require().NotNil(got)
		s.Equal(uuid.Nil, got.Address.ID)
		s.Len(got.Products, 1)
	})

	s.Run("order not found", func() {
		got, err := svc.GetDetailedByID(uuid.New(), owner)
		s.NoError(err)
		s.Nil(got)
	})

	s.Run("order from another user", func() {
		got, err := svc.GetDetailedByID(complete.ID, uuid.New())
		s.NoError(err)
		s.Nil(got, "must not expose orders from other users")
	})
}
