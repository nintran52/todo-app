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
	GetAllItem(userID uuid.UUID, paging *clients.Paging) ([]domain.Item, error)
	GetItemByID(id, userID uuid.UUID) (domain.Item, error)
	UpdateItem(id, userID uuid.UUID, item *domain.ItemUpdate) error
	DeleteItem(id, userID uuid.UUID) error
}

type itemHandler struct {
	itemService ItemService
}

func NewItemHandler(apiVersion *gin.RouterGroup, svc ItemService, middlewareAuth func(c *gin.Context), middlewareRateLimit func(c *gin.Context)) {
	itemHandler := &itemHandler{
		itemService: svc,
	}

	items := apiVersion.Group("/items", middlewareAuth)
	items.POST("", itemHandler.CreateItemHandler)
	items.GET("", middlewareRateLimit, itemHandler.GetAllItemHandler)
	items.GET("/:id", itemHandler.GetItemHandler)
	items.PATCH("/:id", itemHandler.UpdateItemHandler)
	items.DELETE("/:id", itemHandler.DeleteItemHandler)
}

// CreateItemHandler handles the creation of a new item.
//
// @Summary      Create a new item
// @Description  This endpoint allows authenticated users to create an item.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        item  body      domain.ItemCreation  true  "Item creation payload"
// @Success      200   {object}  clients.SuccessRes   "Item successfully created"
// @Failure      400   {object}  clients.AppError     "Bad Request"
// @Failure      401   {object}  clients.AppError     "Unauthorized"
// @Failure      500   {object}  clients.AppError     "Internal Server Error"
// @Router       /items [post]
func (h *itemHandler) CreateItemHandler(c *gin.Context) {
	var item domain.ItemCreation

	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	requester := c.MustGet(clients.CurrentUser).(clients.Requester)
	item.UserID = requester.GetUserID()
	if err := h.itemService.CreateItem(&item); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(item.ID))
}

// GetAllItemHandler retrieves all items.
//
// @Summary      Get all items
// @Description  This endpoint retrieves a list of all items.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Success      200  {object}  clients.SuccessRes  "List of items retrieved successfully"
// @Failure      500  {object}  clients.AppError    "Internal Server Error"
// @Router       /items [get]
func (h *itemHandler) GetAllItemHandler(c *gin.Context) {
	var paging clients.Paging
	if err := c.ShouldBind(&paging); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}
	paging.Process()

	requester := c.MustGet(clients.CurrentUser).(clients.Requester)

	items, err := h.itemService.GetAllItem(requester.GetUserID(), &paging)
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	c.JSON(http.StatusOK, clients.NewSuccessResponse(items, paging, nil))
}

// GetItemHandler retrieves an item by its ID.
//
// @Summary      Get an item by ID
// @Description  This endpoint retrieves a single item by its unique identifier.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true  "Item ID"
// @Success      200  {object}  clients.SuccessRes     "Item retrieved successfully"
// @Failure      400  {object}  clients.AppError       "Invalid ID format or bad request"
// @Failure      404  {object}  clients.AppError       "Item not found"
// @Failure      500  {object}  clients.AppError       "Internal Server Error"
// @Router       /items/{id} [get]
func (h *itemHandler) GetItemHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	requester := c.MustGet(clients.CurrentUser).(clients.Requester)

	item, err := h.itemService.GetItemByID(id, requester.GetUserID())
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(item))
}

// UpdateItemHandler updates an existing item.
//
// @Summary      Update an item
// @Description  This endpoint allows updating the properties of an existing item by its ID.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "Item ID"
// @Param        item  body      domain.ItemUpdate      true  "Item update payload"
// @Success      200   {object}  clients.SuccessRes     "Item updated successfully"
// @Failure      400   {object}  clients.AppError       "Invalid input or bad request"
// @Failure      404   {object}  clients.AppError       "Item not found"
// @Failure      500   {object}  clients.AppError       "Internal Server Error"
// @Router       /items/{id} [put]
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

	requester := c.MustGet(clients.CurrentUser).(clients.Requester)

	if err := h.itemService.UpdateItem(id, requester.GetUserID(), &item); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(true))
}

// DeleteItemHandler deletes an item by its ID.
//
// @Summary      Delete an item
// @Description  This endpoint deletes an item identified by its unique ID.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true  "Item ID"
// @Success      200  {object}  clients.SuccessRes     "Item deleted successfully"
// @Failure      400  {object}  clients.AppError       "Invalid ID format or bad request"
// @Failure      404  {object}  clients.AppError       "Item not found"
// @Failure      500  {object}  clients.AppError       "Internal Server Error"
// @Router       /items/{id} [delete]
func (h *itemHandler) DeleteItemHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	requester := c.MustGet(clients.CurrentUser).(clients.Requester)

	if err := h.itemService.DeleteItem(id, requester.GetUserID()); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(true))
}
