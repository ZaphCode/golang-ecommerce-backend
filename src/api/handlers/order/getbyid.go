package order

import (
	"net/http"

	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Get detailed order handler
// @Summary      Get auth user detailed order
// @Description  Get detailed order by id
// @Tags         order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.DetailedOrderRespOKDTO
// @Param        id   path      string  true  "Order ID"
// @Failure      400  {object}  dtos.DetailRespErrDTO
// @Failure      404  {object}  dtos.DetailRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /order/get/{id} [get]
func (h *OrderHandler) GetOrderById(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, http.StatusBadRequest, "bad request", "invalid order id")
	}

	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return h.RespErr(c, http.StatusInternalServerError, "internal server error", "something went wrong")
	}

	do, err := h.ordSvc.GetDetailedByID(uid, ud.ID)

	if err != nil {
		return h.RespErr(c, http.StatusInternalServerError, "internal server error", "error getting user order")
	}

	if do == nil {
		return h.RespErr(c, http.StatusNotFound, "order not found")
	}

	return h.RespOK(c, http.StatusOK, "order", do)
}
