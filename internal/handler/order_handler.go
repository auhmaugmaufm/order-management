package handler

import (
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/auhmaugmaufm/event-driven-order/internal/service"
	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderHandler struct {
	service *service.OrderService
	cfg     *config.Config
}

func NewOrderHandler(svc *service.OrderService, cfg *config.Config) *OrderHandler {
	return &OrderHandler{service: svc, cfg: cfg}
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var req dto.OrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid request body",
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	err := h.service.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to create order",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *OrderHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid order ID format",
		})
	}

	order, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: "order not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   order,
		Status: fiber.StatusOK,
	})
}

func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	orders, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to fetch orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   orders,
		Status: fiber.StatusOK,
	})
}
