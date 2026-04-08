package handler

import (
	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
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
			Message: err.Error(),
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

	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			OrderID:   item.OrderID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	data := &dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   data,
		Status: fiber.StatusOK,
	})
}

func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	var req dto.PaginationRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: err.Error(),
		})
	}
	req.SetDefaults()
	pagination := &domain.Pagination{
		Limit: req.Limit,
		Page:  req.Page,
	}
	orders, total, err := h.service.GetAll(c.Context(), pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to fetch orders",
		})
	}
	totalPage := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)
	return c.Status(fiber.StatusOK).JSON(dto.PaginationResponse{
		Data:        orders,
		TotalItems:  total,
		TotalPages:  totalPage,
		CurrentPage: pagination.Page,
		Status:      fiber.StatusOK,
	})
}
