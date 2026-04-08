package handler

import (
	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/auhmaugmaufm/event-driven-order/internal/service"
	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StockHandler struct {
	service *service.StockService
	cfg     *config.Config
}

func NewStockHandler(svc *service.StockService, cfg *config.Config) *StockHandler {
	return &StockHandler{service: svc, cfg: cfg}
}

// TODO: Stock Adjustment by create stock movement

// func (h *StockHandler) IncreaseStock(c *fiber.Ctx) error {
// 	var req dto.UpdateStockReq
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error:   "bad_request",
// 			Message: err.Error(),
// 		})
// 	}

// 	err := h.service.IncreaseStock(req.ProductID, req.Quantity)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
// 			Error:   "not_found",
// 			Message: err.Error(),
// 		})
// 	}
// 	return c.SendStatus(fiber.StatusOK)
// }

// func (h *StockHandler) DecreaseStock(c *fiber.Ctx) error {
// 	var req dto.UpdateStockReq
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error:   "bad_request",
// 			Message: err.Error(),
// 		})
// 	}

// 	err := h.service.DecreaseStock(req.ProductID, req.Quantity)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
// 			Error:   "not_found",
// 			Message: err.Error(),
// 		})
// 	}
// 	return c.SendStatus(fiber.StatusOK)
// }

func (h *StockHandler) GetProductStock(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid product id",
		})
	}

	res, err := h.service.GetProductStock(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}

	response := &dto.StockResponse{
		ID:          res.ID,
		ProductID:   res.ProductID,
		ProductName: res.Product.ProductName,
		Quantity:    res.Quantity,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   response,
		Status: fiber.StatusOK,
	})

}

func (h *StockHandler) GetAllProductStocks(c *fiber.Ctx) error {
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
	res, total, err := h.service.GetAll(c.Context(), pagination)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}

	response := make([]dto.StockResponse, 0, len(res))

	for _, stock := range res {
		productName := ""
		if stock.Product != nil {
			productName = stock.Product.ProductName
		}
		response = append(response, dto.StockResponse{
			ID:          stock.ID,
			ProductID:   stock.ProductID,
			ProductName: productName,
			Quantity:    stock.Quantity,
			CreatedAt:   stock.CreatedAt,
			UpdatedAt:   stock.UpdatedAt,
		})
	}

	totalPage := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)
	return c.Status(fiber.StatusOK).JSON(dto.PaginationResponse{
		Data:        response,
		TotalItems:  total,
		TotalPages:  totalPage,
		CurrentPage: pagination.Page,
		Status:      fiber.StatusOK,
	})

}
