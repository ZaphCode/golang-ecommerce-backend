package dtos

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

type NewOrderDTO struct {
	PaymentID string                `json:"payment_id" validate:"required" example:"pm_1NKPiEG8UXDxPRbaEDuh6BrU"`
	Products  []domain.OrderProduct `json:"products" validate:"required"`
	AddressID uuid.UUID             `json:"address_id" validate:"required" example:"8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"`
}

func (dto NewOrderDTO) AdaptToOrder(price int64, usrid uuid.UUID) domain.Order {
	return domain.Order{
		UserID:    usrid,
		AddressID: dto.AddressID,
		PaymentID: dto.PaymentID,
		Amount:    price,
		Products:  dto.Products,
		Status:    utils.StatusPending,
		Paid:      false,
	}
}

type OrderDTO struct {
	NewOrderDTO
	ID        uuid.UUID `json:"id" example:"8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"`
	CreatedAt int64     `json:"created_at" example:"1674405183"`
	Amount    int64     `json:"amount" example:"14500"`
	Status    string    `json:"status" example:"pending"`
	Paid      bool      `json:"paid" example:"true"`
	UpdatedAt int64     `json:"updated_at" example:"1674405181"`
}

type DetailedOrderDTO struct {
	OrderDTO
	Address  AddressDTO                `json:"address"`
	Products []DetailedOrderProductDTO `json:"products"`
}

type DetailedOrderProductDTO struct {
	ProductDTO
	Quantity uint `json:"quantity" example:"2"`
}
