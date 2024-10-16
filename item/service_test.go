package item_test

import (
	"errors"
	"testing"
	"todo-app/domain"
	service "todo-app/item"
	"todo-app/item/mocks"
	"todo-app/pkg/clients"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateItem(t *testing.T) {
	// Create a mock item repository
	mockItemRepo := new(mocks.ItemRepo)
	itemService := service.NewItemService(mockItemRepo)

	t.Run("success", func(t *testing.T) {
		// Define valid item creation input
		item := &domain.ItemCreation{
			Title:       "New Item",
			Description: "New Item Description",
			UserID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		}

		// Setup mock expectation
		mockItemRepo.On("Save", mock.Anything).Return(nil).Once()

		// Call the service method
		err := itemService.CreateItem(item)

		// Assertions
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, item.ID) // Ensure the ID is generated

		// Verify that all expectations were met
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		// Define invalid item input (e.g., missing title)
		item := &domain.ItemCreation{
			Description: "Invalid Item",
			UserID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		}

		// Call the service method
		err := itemService.CreateItem(item)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title can not be null")
	})

	t.Run("repository error", func(t *testing.T) {
		// Define valid item input
		item := &domain.ItemCreation{
			Title:       "Another Item",
			Description: "Another Item Description",
			UserID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		}

		// Setup mock expectation to simulate a save failure
		mockItemRepo.On("Save", mock.Anything).Return(errors.New("cannot create entity")).Once()

		// Call the service method
		err := itemService.CreateItem(item)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot create entity")

		// Verify that all expectations were met
		mockItemRepo.AssertExpectations(t)
	})
}

// TestGetAllItem tests the GetAllItem method of itemService.
func TestGetAllItem(t *testing.T) {
	// Define mock data
	mockItems := []domain.Item{
		{
			ID:          uuid.New(),
			Title:       "Item 1",
			Description: "Description 1",
			UserID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		},
		{
			ID:          uuid.New(),
			Title:       "Item 2",
			Description: "Description 2",
			UserID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		},
	}

	// Create a mock item repository
	mockItemRepo := new(mocks.ItemRepo)
	itemService := service.NewItemService(mockItemRepo) // Create the service with the mock repo

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	paging := &clients.Paging{Limit: 2, Page: 1}

	t.Run("success", func(t *testing.T) {
		// Setup mock expectation for success
		mockItemRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*clients.Paging")).
			Return(mockItems, nil).Once()

		// Call the service method
		result, err := itemService.GetAllItem(userID, paging)

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Item 1", result[0].Title)
		assert.Equal(t, "Item 2", result[1].Title)

		// Verify that all expectations were met
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Simulate a repository error
		mockErr := errors.New("repository error")
		mockItemRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*clients.Paging")).
			Return(nil, mockErr).Once()

		// Call the service method
		result, err := itemService.GetAllItem(userID, paging)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify that all expectations were met
		mockItemRepo.AssertExpectations(t)
	})
}

func TestGetItemByID(t *testing.T) {
	// Create a mock item repository
	mockItemRepo := new(mocks.ItemRepo)
	itemService := service.NewItemService(mockItemRepo)

	// Define mock data
	mockID := uuid.New()
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	mockItem := domain.Item{
		ID:          mockID,
		Title:       "Test Item",
		Description: "Test Description",
		UserID:      userID,
	}

	t.Run("success", func(t *testing.T) {
		// Setup mock expectation for a successful fetch
		mockItemRepo.On("GetItem", mock.Anything).
			Return(mockItem, nil).Once()

		// Call the service method
		result, err := itemService.GetItemByID(mockID, userID)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, mockItem, result)

		// Verify expectations
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("error - item not found", func(t *testing.T) {
		// Simulate item not found scenario
		mockItemRepo.On("GetItem", mock.Anything).
			Return(domain.Item{}, errors.New("item not found")).Once()

		// Call the service method
		result, err := itemService.GetItemByID(mockID, userID)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, domain.Item{}, result)

		// Verify expectations
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		// Simulate a repository error
		mockItemRepo.On("GetItem", mock.Anything).
			Return(domain.Item{}, errors.New("cannot get entity")).Once()

		// Call the service method
		_, err := itemService.GetItemByID(mockID, userID)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot get entity")

		// Verify expectations
		mockItemRepo.AssertExpectations(t)
	})
}

func TestUpdateItem(t *testing.T) {
	// Create a mock item repository
	mockItemRepo := new(mocks.ItemRepo)
	itemService := service.NewItemService(mockItemRepo)

	// Define mock data
	mockID := uuid.New()
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	updateData := &domain.ItemUpdate{
		Title:       ptrToString("Updated Title"),
		Description: ptrToString("Updated Description"),
	}

	t.Run("success", func(t *testing.T) {
		// Setup mock expectation for successful update
		mockItemRepo.On("Update", mock.Anything, updateData).Return(nil).Once()

		// Call the service method
		err := itemService.UpdateItem(mockID, userID, updateData)

		// Assertions
		assert.NoError(t, err)

		// Verify expectations
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		// Simulate repository error
		mockItemRepo.On("Update", mock.Anything, updateData).
			Return(errors.New("cannot update entity")).Once()

		// Call the service method
		err := itemService.UpdateItem(mockID, userID, updateData)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot update entity")

		// Verify expectations
		mockItemRepo.AssertExpectations(t)
	})
}

func TestDeleteItem(t *testing.T) {
	// Create a mock item repository
	mockItemRepo := new(mocks.ItemRepo)
	itemService := service.NewItemService(mockItemRepo)

	// Define mock data
	mockID := uuid.New()
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("success", func(t *testing.T) {
		// Setup mock expectation for successful deletion
		mockItemRepo.On("Delete", mock.Anything).Return(nil).Once()

		// Call the service method
		err := itemService.DeleteItem(mockID, userID)

		// Assertions
		assert.NoError(t, err)

		// Verify expectations
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		// Simulate a repository error
		mockItemRepo.On("Delete", mock.Anything).
			Return(errors.New("cannot delete entity")).Once()

		// Call the service method
		err := itemService.DeleteItem(mockID, userID)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete entity")

		// Verify expectations
		mockItemRepo.AssertExpectations(t)
	})
}

// Helper function to return a pointer to a string
func ptrToString(s string) *string {
	return &s
}
