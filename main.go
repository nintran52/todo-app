package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Status int

const (
	Deleted Status = iota
	Active
	Done
)

type Item struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (Item) TableName() string { return "items" }

type ItemCreation struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (ItemCreation) TableName() string { return Item{}.TableName() }

type ItemUpdate struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Status      *Status   `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ItemUpdate) TableName() string { return Item{}.TableName() }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("CONNECTION_STRING")), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	api := r.Group("v1")
	{
		items := api.Group("items")
		items.POST("", CreateItem(db))
		items.GET("", GetAllItems(db))
		items.GET("/:id", GetItem(db))
		items.PATCH("/:id", UpdateItem(db))
		items.DELETE("/:id", DeleteItem(db))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CreateItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		itemData := ItemCreation{}

		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		itemData.ID = uuid.New()
		if err := db.Create(&itemData).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": itemData.ID,
		})
	}
}

func GetAllItems(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		items := []Item{}

		if err := db.Find(&items).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": items,
		})
	}
}

func GetItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		itemData := Item{}
		if err := db.Where("id = ?", id).Find(&itemData).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": itemData,
		})
	}
}

func UpdateItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		itemData := ItemUpdate{}
		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		itemData.UpdatedAt = time.Now()
		if err := db.Where("id = ?", id).Updates(&itemData).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func DeleteItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		// // Hard delete
		// if err := db.Table(Item{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err.Error(),
		// 	})

		// 	return
		// }

		// Soft delete
		deletedStatus := Deleted
		if err := db.Where("id = ?", id).Updates(&ItemUpdate{Status: &deletedStatus}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}
