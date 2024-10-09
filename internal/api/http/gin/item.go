package gin

import (
	"net/http"
	"todo-app/domain"
	"todo-app/pkg/clients"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ItemService interface {
	CreateItem(item *domain.ItemCreation) error
	GetAllItem() ([]domain.Item, error)
	GetItemByID(id uuid.UUID) (domain.Item, error)
	UpdateItem(id uuid.UUID, item *domain.ItemUpdate) error
	DeleteItem(id uuid.UUID) error
}

type itemHandler struct {
	itemService ItemService
}

func NewItemHandler(apiVersion *gin.RouterGroup, svc ItemService) {
	itemHandler := &itemHandler{
		itemService: svc,
	}

	items := apiVersion.Group("/items")
	items.POST("", itemHandler.CreateItemHandler)
	items.GET("", itemHandler.GetAllItemHandler)
	items.GET("/:id", itemHandler.GetItemHandler)
	items.PATCH("/:id", itemHandler.UpdateItemHandler)
	items.DELETE("/:id", itemHandler.DeleteItemHandler)
}

func (h *itemHandler) CreateItemHandler(c *gin.Context) {
	var item domain.ItemCreation

	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	if err := h.itemService.CreateItem(&item); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(item.ID))
}

func (h *itemHandler) GetAllItemHandler(c *gin.Context) {
	items, err := h.itemService.GetAllItem()
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(items))
}

func (h *itemHandler) GetItemHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	item, err := h.itemService.GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(item))
}

func (h *itemHandler) UpdateItemHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	item := domain.ItemUpdate{}
	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	if err := h.itemService.UpdateItem(id, &item); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(true))
}

func (h *itemHandler) DeleteItemHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	if err := h.itemService.DeleteItem(id); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(true))
}
