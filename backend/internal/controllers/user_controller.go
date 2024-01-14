package controllers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/kiloMIA/mediscan/backend/internal/models"
)

type UserController struct {
	DB *pgxpool.Pool
	lg *logrus.Logger
}

func NewUserController(db *pgxpool.Pool, lg *logrus.Logger) *UserController {
	return &UserController{
		DB: db,
		lg: lg,
	}
}


func (userc *UserController) CreateUser(ctx context.Context, user *models.User) error {
	userc.lg.Debugln("User Creation at controller level")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		userc.lg.Errorf("user controller - CreateUser - password generation - %v", err)
		return err
	}

	_, err = userc.DB.Exec(ctx,
		"INSERT INTO users (name, surname, email, password) VALUES ($1, $2, $3, $4)",
		user.Name, user.Surname, user.Email, string(hashedPassword))

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") && strings.Contains(err.Error(), "Email") {
			return fmt.Errorf("email %s already exists", user.Email)
		}
		userc.lg.Errorf("user controller - CreateUser - db exec - %v", err)
		return err
	}

	return nil
}

func (userc *UserController) Authenticate(ctx context.Context, email, password string) (models.User, error) {
	userc.lg.Debugln("User Authentication at controller level")
	var user models.User
	err := userc.DB.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		userc.lg.Errorf("user controller - Authenticate - db exec - %v", err)
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		userc.lg.Errorf("user controller - Authenticate - hash and password comparison - %v", err)
		return models.User{}, errors.New("incorrect password")
	}

	return user, nil
}

func (userc *UserController) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	userc.lg.Debugln("Getting user by ID at controller level")

	var user models.User
	err := userc.DB.QueryRow(ctx, "SELECT id, name, surname,email, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password)

	if err != nil {
		userc.lg.Errorf("user controller - GetUserByID - db exec - %v", err)
		return nil, err
	}

	return &user, nil
}