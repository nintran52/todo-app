package main

import (
	"log"
	"net/http"
	"os"

	pgRepo "todo-app/internal/repository/postgres"
	"todo-app/item"
	"todo-app/pkg/tokenprovider/jwt"
	"todo-app/pkg/util"
	"todo-app/user"

	restApi "todo-app/internal/api/http/gin"
	"todo-app/internal/api/http/gin/middleware"

	"github.com/gin-gonic/gin"
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
	r.Use(middleware.Recover())

	apiVersion := r.Group("v1")

	itemRepo := pgRepo.NewItemRepo(db)
	itemService := item.NewItemService(itemRepo)

	userRepo := pgRepo.NewUserRepo(db)
	hasher := util.NewMd5Hash()
	tokenProvider := jwt.NewJWTProvider()
	userService := user.NewUserService(userRepo, hasher, tokenProvider, 60*60*24*30)

	restApi.NewItemHandler(apiVersion, itemService)
	restApi.NewUserHandler(apiVersion, userService)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
