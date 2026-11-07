package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

type orderService struct {
	ordRepo  domain.OrderRepository
	addrRepo domain.AddressRepository
	prodRepo domain.ProductRepository
}

func NewOrderService(
	ordRepo domain.OrderRepository,
	addrRepo domain.AddressRepository,
	prodRepo domain.ProductRepository,
) domain.OrderService {
	return &orderService{
		ordRepo:  ordRepo,
		addrRepo: addrRepo,
		prodRepo: prodRepo,
	}
}

func (s *orderService) Create(ord *domain.Order) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating uuid: %s", err)
	}

	if addr, err := s.addrRepo.FindByID(ord.AddressID); err != nil || addr == nil {
		return fmt.Errorf("invalid address id")
	}

	ord.ID = ID
	ord.Status = utils.StatusProcessing
	ord.CreatedAt = time.Now().Unix()
	ord.UpdatedAt = time.Now().Unix()

	if err := s.ordRepo.Save(ord); err != nil {
		return fmt.Errorf("error saving order: %w", err)
	}

	return nil
}

// GetDetailedByID returns the order with its address and products loaded.
// It returns (nil, nil) when the order doesn't exist or doesn't belong to userID,
// so callers can't tell other users' orders apart from missing ones.
func (s *orderService) GetDetailedByID(ID, userID uuid.UUID) (*domain.DetailedOrder, error) {
	ord, err := s.ordRepo.FindByID(ID)

	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	// Repositories return (nil, nil) when the resource doesn't exist.
	if ord == nil || ord.UserID != userID {
		return nil, nil
	}

	detailedOrder := domain.DetailedOrder{Order: *ord}

	addr, err := s.addrRepo.FindByID(ord.AddressID)

	if err != nil {
		return nil, fmt.Errorf("error getting address: %w", err)
	}

	// The address may have been deleted after the order was placed:
	// keep the order and leave the address empty.
	if addr != nil {
		detailedOrder.Address = *addr
	}

	for _, op := range ord.Products {
		prod, err := s.prodRepo.FindByID(op.ID)

		if err != nil {
			return nil, fmt.Errorf("error getting product: %w", err)
		}

		// Skip products that no longer exist instead of failing the whole order.
		if prod != nil {
			detailedOrder.Products = append(detailedOrder.Products, domain.DetailedOrderProduct{
				Product:  *prod,
				Quantity: op.Quantity,
			})
		}
	}

	return &detailedOrder, nil
}

func (s *orderService) GetAll() ([]domain.Order, error) {
	return s.ordRepo.Find()
}

func (s *orderService) UpdateStatus(ID uuid.UUID, status string) error {
	return s.ordRepo.UpdateField(ID, "Status", status)
}

func (s *orderService) SetPaidStatus(ID uuid.UUID, paid bool) error {
	return s.ordRepo.UpdateField(ID, "Paid", paid)
}

func (s *orderService) GetAllByUserID(ID uuid.UUID) ([]domain.Order, error) {
	return s.ordRepo.FindWhere("UserID", "==", ID)
}

func (s *orderService) Delete(ID uuid.UUID) error {
	return s.ordRepo.Remove(ID)
}
