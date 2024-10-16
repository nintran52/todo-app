package postgres_test

import (
	"errors"
	"testing"
	"time"
	"todo-app/domain"
	"todo-app/internal/repository/postgres"
	"todo-app/item"
	"todo-app/pkg/clients"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup function for creating the in-memory SQLite database.
func setupTestDB() (*gorm.DB, item.ItemRepo, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	if err := db.AutoMigrate(&domain.Item{}); err != nil {
		return nil, nil, err
	}
	return db, postgres.NewItemRepo(db), nil
}

// Helper function to insert a mock item into the database.
func insertMockItem(db *gorm.DB, title, description string, userID uuid.UUID) domain.Item {
	item := domain.Item{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		UserID:      userID,
		Status:      1,
	}

	db.Create(&item)

	return item
}

func TestSaveItem(t *testing.T) {
	db, repo, err := setupTestDB()
	require.NoError(t, err)

	item := &domain.ItemCreation{
		ID:          uuid.New(),
		Title:       "Test Item",
		Description: "Test Description",
		UserID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
	}

	err = repo.Save(item)
	assert.NoError(t, err)

	var savedItem domain.Item
	err = db.First(&savedItem, "id = ?", item.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, item.Title, savedItem.Title)
}

func TestGetAllItems(t *testing.T) {
	db, repo, err := setupTestDB()
	require.NoError(t, err)

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	insertMockItem(db, "Item 1", "Description 1", userID)
	insertMockItem(db, "Item 2", "Description 2", userID)

	filter := map[string]any{"user_id": userID}
	paging := &clients.Paging{Limit: 2, Page: 1}

	result, err := repo.GetAll(filter, paging)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Item 1", result[0].Title)
	assert.Equal(t, "Item 2", result[1].Title)
}

func TestGetItem_Success(t *testing.T) {
	db, repo, err := setupTestDB()
	require.NoError(t, err)

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	item := insertMockItem(db, "Test Item", "Test Description", userID)

	result, err := repo.GetItem(map[string]any{"id": item.ID})
	assert.NoError(t, err)
	assert.Equal(t, item.ID, result.ID)
	assert.Equal(t, item.Title, result.Title)
}

func TestGetItem_NotFound(t *testing.T) {
	_, repo, err := setupTestDB()
	require.NoError(t, err)

	_, err = repo.GetItem(map[string]any{"id": uuid.New()})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, clients.ErrRecordNotFound))
}

func TestUpdateItem(t *testing.T) {
	db, repo, err := setupTestDB()
	require.NoError(t, err)

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	item := insertMockItem(db, "Old Title", "Old Description", userID)

	title := "New Title"
	description := "New Description"
	updateData := domain.ItemUpdate{
		Title:       &title,
		Description: &description,
		UpdatedAt:   time.Now(),
	}

	err = repo.Update(map[string]any{"id": item.ID}, &updateData)
	assert.NoError(t, err)

	var updatedItem domain.Item
	err = db.First(&updatedItem, "id = ?", item.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "New Title", updatedItem.Title)
	assert.Equal(t, "New Description", updatedItem.Description)
}

func TestDeleteItem(t *testing.T) {
	db, repo, err := setupTestDB()
	require.NoError(t, err)

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	item := insertMockItem(db, "Test Item", "Test Description", userID)

	err = repo.Delete(map[string]any{"id": item.ID})
	assert.NoError(t, err)

	var deletedItem domain.Item
	err = db.First(&deletedItem, "id = ?", item.ID).Error
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
