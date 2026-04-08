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

type ProductHandler struct {
	service *service.ProductService
	cfg     *config.Config
}

func NewProductHandler(svc *service.ProductService, cfg *config.Config) *ProductHandler {
	return &ProductHandler{service: svc, cfg: cfg}
}

var validate = validator.New()

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req dto.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	if err := h.service.Create(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to create product",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "createad successfully",
	})
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid product id",
		})
	}

	res, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}

	response := &dto.ProductResponse{
		ID:           res.ID,
		ProductName:  res.ProductName,
		ProductPrice: res.ProductPrice,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   response,
		Status: fiber.StatusOK,
	})
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
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
	products, total, err := h.service.GetAll(c.Context(), pagination)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}

	response := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			ID:           product.ID,
			ProductName:  product.ProductName,
			ProductPrice: product.ProductPrice,
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
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
