package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todo-app/docs"
	restApi "todo-app/internal/api/http/gin"
	"todo-app/internal/api/http/gin/middleware"
	pgRepo "todo-app/internal/repository/postgres"
	"todo-app/item"
	"todo-app/pkg/memcache"
	"todo-app/pkg/tokenprovider/jwt"
	"todo-app/pkg/util"
	"todo-app/user"
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
	docs.SwaggerInfo.BasePath = "/v1"

	itemRepo := pgRepo.NewItemRepo(db)
	itemService := item.NewItemService(itemRepo)

	userRepo := pgRepo.NewUserRepo(db)
	hasher := util.NewMd5Hash()
	tokenProvider := jwt.NewJWTProvider(os.Getenv("SECRET_KEY"))
	tokenExpire := 60 * 60 * 24 * 30
	userService := user.NewUserService(userRepo, hasher, tokenProvider, tokenExpire)

	authCache := memcache.NewUserCaching(memcache.NewRedisCache(), userRepo)
	middlewareAuth := middleware.RequiredAuth(tokenProvider, authCache)

	limiterRate := limiter.Rate{
		Period: 5 * time.Second,
		Limit:  3,
	}
	store := memory.NewStore()
	limiter := limiter.New(store, limiterRate)
	middlewareRateLimit := middleware.RateLimiter(limiter)

	restApi.NewItemHandler(apiVersion, itemService, middlewareAuth, middlewareRateLimit)
	restApi.NewUserHandler(apiVersion, userService)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
