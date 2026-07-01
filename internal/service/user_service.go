package service

import (
	"errors"
	"fmt"
	"time"
	"user-service/internal/model"

	jwtutil "github.com/miank1/ecommerce_backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	GetByID(string) (*model.User, error)
}

type UserService struct {
	Repo      UserRepository
	JWTTTL    int // minutes
	JWTSecret string
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{Repo: r}
}

// Login validates credentials and returns JWT + user (without password)
func (s *UserService) Login(email, password string) (string, *model.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// generate token

	fmt.Println("BEFORE TOKEN GENERATION ")
	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	user.PasswordHash = ""
	return token, user, nil
}

func (s *UserService) Register(name, email, password string) (*model.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email and password are required")
	}
	// check email uniqueness
	existing, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.Repo.Create(u); err != nil {
		return nil, err
	}

	// don't return password hash
	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) Authenticate(email, password string) (*model.User, error) {

	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
