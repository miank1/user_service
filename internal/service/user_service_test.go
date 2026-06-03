package service

import (
	"testing"
	"user-service/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	User *model.User
	Err  error
}

func (m *MockUserRepo) Create(*model.User) error {
	return nil
}

func (m *MockUserRepo) FindByEmail(email string) (*model.User, error) {
	return m.User, m.Err
}

func (m *MockUserRepo) GetByID(id string) (*model.User, error) {
	return m.User, m.Err
}

func TestAuthenticate_Success(t *testing.T) {

	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("admin"),
		bcrypt.DefaultCost,
	)

	mockRepo := &MockUserRepo{
		User: &model.User{
			ID:           "123",
			Email:        "admin@gmail.com",
			PasswordHash: string(hash),
		},
	}

	svc := NewUserService(mockRepo)

	user, err := svc.Authenticate(
		"admin@gmail.com",
		"admin",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user")
	}
}

func TestAuthenticate_InvalidPassword(t *testing.T) {

	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("admin"),
		bcrypt.DefaultCost,
	)

	mockRepo := &MockUserRepo{
		User: &model.User{
			ID:           "123",
			Email:        "admin@gmail.com",
			PasswordHash: string(hash),
		},
	}

	svc := NewUserService(mockRepo)

	_, err := svc.Authenticate(
		"admin@gmail.com",
		"wrong-password",
	)

	if err == nil {
		t.Fatal("expected error")
	}
}
