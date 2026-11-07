package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Add addr repo to user service

type userService struct {
	usrRepo domain.UserRepository
}

func NewUserService(usrRepo domain.UserRepository) domain.UserService {
	return &userService{usrRepo: usrRepo}
}

func (s *userService) Create(user *domain.User) error {
	eusr, err := s.usrRepo.FindByField("Email", user.Email)

	if err != nil {
		return fmt.Errorf("internal server error: %s", err)
	}

	if eusr != nil {
		return fmt.Errorf("email taken")
	}

	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("uuid generation error: %s", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error hashing password: %s", err)
	}

	user.ID = ID
	user.Password = string(hash)
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	if user.Age == 0 {
		user.Age = 18
	}

	if user.Role == "" {
		user.Role = utils.UserRole
	}

	if user.ImageUrl == "" {
		user.ImageUrl = fmt.Sprintf("https://api.dicebear.com/5.x/bottts-neutral/svg?seed=%d", user.CreatedAt)
	}

	if err := s.usrRepo.Save(user); err != nil {
		return err
	}

	user.Password = ""

	return nil
}

func (s *userService) GetAll() ([]domain.User, error) {
	users, err := s.usrRepo.Find()

	if err != nil {
		return nil, err
	}

	for i := range users {
		hidePassword(&users[i])
	}

	return users, nil
}

func (s *userService) GetByID(ID uuid.UUID) (*domain.User, error) {
	user, err := s.usrRepo.FindByID(ID)

	if err != nil {
		return nil, err
	}

	hidePassword(user)

	return user, nil
}

func (s *userService) GetByEmail(email string) (*domain.User, error) {
	user, err := s.usrRepo.FindByField("Email", email)

	if err != nil {
		return nil, err
	}

	hidePassword(user)

	return user, nil
}

func (s *userService) VerifyEmail(ID uuid.UUID) error {
	return s.usrRepo.UpdateField(ID, "VerifiedEmail", true)
}

func (s *userService) GetByCredentials(email, password string) (*domain.User, error) {
	user, err := s.usrRepo.FindByField("Email", email)

	if err != nil {
		return nil, fmt.Errorf("internal server error: %s", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	hidePassword(user)

	return user, nil
}

func (s *userService) UpdatePassword(ID uuid.UUID, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error hasing password %w", err)
	}

	return s.usrRepo.UpdateField(ID, "Password", string(hash))
}

func (s *userService) Update(ID uuid.UUID, uf domain.UpdateFields) error {
	delete(uf, "Email")
	delete(uf, "Model")
	return s.usrRepo.Update(ID, uf)
}

func (s *userService) Delete(ID uuid.UUID) error {
	return s.usrRepo.Remove(ID)
}

// Helper functions

func hidePassword(user *domain.User) {
	if user != nil {
		user.Password = ""
	}
}
