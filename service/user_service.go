package service

import (
	"errors"
	"go-crud/input"
	"go-crud/models"
	"go-crud/repository"

	"golang.org/x/crypto/bcrypt"
)

type ServiceUser interface {
	RegisterUser(inputUser input.UserInput) (*models.User, error)
	LoginUser(inputUser input.UserInput) (*models.User, error)
	GetUserByid(ID int) (*models.User, error)
	IsEmailAvailability(input string) (bool, error)
}

type serviceUser struct {
	repositoryUser repository.UserRepository
}

func NewUserService(repositoryUser repository.UserRepository) *serviceUser {
	return &serviceUser{repositoryUser}
}

func (s *serviceUser) RegisterUser(inputUser input.UserInput) (*models.User, error) {
	user := &models.User{
		Email:    inputUser.Email,
		Password: inputUser.Password,
		Role:     0,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(inputUser.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.Password = string(passwordHash)

	createUser, err := s.repositoryUser.CreateUser(user)
	if err != nil {
		return createUser, err
	}
	return createUser, nil
}

func (s *serviceUser) LoginUser(inputUser input.UserInput) (*models.User, error) {
	email := inputUser.Email
	password := inputUser.Password

	checkUser, err := s.repositoryUser.FindByEmail(email)
	if err != nil {
		return checkUser, err
	}
	if checkUser.ID == 0 {
		return checkUser, errors.New("user not found that email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(password))
	if err != nil {
		return nil, err // Password tidak cocok
	}

	return checkUser, nil
}

func (s *serviceUser) IsEmailAvailability(input string) (bool, error) {
	user, err := s.repositoryUser.FindByEmail(input)
	if err != nil {
		return false, err
	}

	// Jika user tidak ditemukan (nil), email tersedia
	if user == nil {
		return true, nil
	}

	// Jika user ditemukan, email sudah digunakan
	return false, nil
}

func (s *serviceUser) GetUserByid(ID int) (*models.User, error) {
	user, err := s.repositoryUser.FindById(ID)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("user Not Found With That ID")
	}

	return user, nil
}
