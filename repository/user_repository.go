package repository

import (
	"database/sql"
	"go-crud/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindById(ID int) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {

	query := "INSERT INTO users (Email, Password, Role) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, user.Email, user.Password, user.Role)
	if err != nil {
		return user, err
	}

	userID, _ := result.LastInsertId()
	user.ID = int(userID)

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	query := "SELECT id, email, password, role FROM users WHERE email = ?"

	row := r.db.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err // Other errors
	}

	return user, nil

}

func (r *userRepository) FindById(ID int) (*models.User, error) {
	query := "SELECT id, email, password, role FROM users WHERE id = ?"

	row := r.db.QueryRow(query, ID)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err // Other errors
	}

	return user, nil
}
