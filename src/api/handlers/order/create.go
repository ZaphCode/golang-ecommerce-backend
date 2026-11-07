package order

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
)

// * Create new order handler
// @Summary      Create new order
// @Description  Create order
// @Tags         order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_data  body dtos.NewOrderDTO true "order data"
// @Success      200  {object}  dtos.OrderRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /order/new [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	usrData, ok1 := c.Locals("user-data").(*auth.Claims)
	cusID, ok2 := c.Locals("customer-id").(string)

	if !ok1 || !ok2 {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	body := dtos.NewOrderDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	price, err := h.prodSvc.CalculateTotalPrice(body.Products)

	if err != nil {
		return h.RespErr(c, 400, "some product are invalid", err.Error())
	}

	order := body.AdaptToOrder(price, usrData.ID)

	if err := h.ordSvc.Create(&order); err != nil {
		return h.RespErr(c, 500, "error creating order", err.Error())
	}

	err = h.pmSvc.MakePayment(cusID, body.PaymentID, price)

	if err != nil {
		return h.RespErr(c, 500, "error making the payment", err.Error())
	}

	go func() {
		if err := h.ordSvc.SetPaidStatus(order.ID, true); err != nil {
			utils.PrintColor("red", "Error updating order status")
		}
	}()

	return h.RespOK(c, 200, "order created", order)
}
