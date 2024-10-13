package middleware

import (
	"errors"
	"strings"
	"todo-app/domain"
	"todo-app/pkg/clients"
	"todo-app/pkg/tokenprovider"

	"github.com/gin-gonic/gin"
)

type AuthenRepo interface {
	GetUser(conditions map[string]interface{}) (*domain.User, error)
}

func RequiredAuth(tokenProvider tokenprovider.Provider, userRepo AuthenRepo) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		if err != nil {
			panic(err)
		}

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(err)
		}

		user, err := userRepo.GetUser(map[string]interface{}{"id": payload.UserID()})
		if err != nil {
			panic(err)
		}

		if user.Status == 0 {
			panic(clients.ErrNoPermission(errors.New("user has been deleted or banned")))
		}

		c.Set(clients.CurrentUser, user)
		c.Next()
	}
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	//"Authorization" : "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(nil)
	}

	return parts[1], nil
}

func ErrWrongAuthHeader(err error) *clients.AppError {
	return clients.NewCustomError(
		err,
		"wrong authen header",
		"ErrWrongAuthHeader",
	)
}
