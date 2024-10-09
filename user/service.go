package user

import (
	"errors"
	"todo-app/domain"
	"todo-app/pkg/clients"
	"todo-app/pkg/tokenprovider"
	"todo-app/pkg/util"

	"github.com/google/uuid"
)

type UserRepo interface {
	Save(user *domain.UserCreate) error
	GetUser(conditions map[string]any) (*domain.User, error)
}

type Hasher interface {
	Hash(data string) string
}

type userService struct {
	userRepo      UserRepo
	hasher        Hasher
	tokenProvider tokenprovider.Provider
	expiry        int
}

func NewUserService(repo UserRepo, hasher Hasher, tokenProvider tokenprovider.Provider, expiry int) *userService {
	return &userService{
		userRepo:      repo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		expiry:        expiry,
	}
}

func (s *userService) Register(data *domain.UserCreate) error {
	if err := data.Validate(); err != nil {
		return clients.ErrInvalidRequest(err)
	}

	user, err := s.userRepo.GetUser(map[string]any{"email": data.Email})
	if err != nil {
		if !errors.Is(err, clients.ErrRecordNotFound) {
			return err
		}
	}

	if user != nil {
		return domain.ErrEmailExisted
	}

	salt := util.GenSalt(50)

	data.ID = uuid.New()
	data.Password = s.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = 1

	if err := s.userRepo.Save(data); err != nil {
		return clients.ErrCannotCreateEntity(data.TableName(), err)
	}

	return nil
}

func (s *userService) Login(data *domain.UserLogin) (tokenprovider.Token, error) {
	user, err := s.userRepo.GetUser(map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, domain.ErrEmailOrPasswordInvalid
	}

	passHashed := s.hasher.Hash(data.Password + user.Salt)

	if user.Password != passHashed {
		return nil, domain.ErrEmailOrPasswordInvalid
	}

	payload := &clients.TokenPayload{
		UID:   user.ID,
		URole: user.Role.String(),
	}

	accessToken, err := s.tokenProvider.Generate(payload, s.expiry)
	if err != nil {
		return nil, clients.ErrInternal(err)
	}

	return accessToken, nil
}
