package postgres

import (
	"time"

	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	model := &UserModel{
		ID: uuid.New().String(),
		Login: user.Login,
		PasswordHash: user.Password,
		CreatedAt: time.Now(),
	}
	err := r.db.Create(model).Error
	if err == nil{
		user.ID = model.ID
	}
	return err
}

func (r *userRepository) FindByLogin(login string) (*domain.User, error) {
	var model UserModel
	if err := r.db.Where("login = ?", login).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:       model.ID,
		Login:    model.Login,
		Password: model.PasswordHash,
	}, nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var model UserModel
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:       model.ID,
		Login:    model.Login,
		Password: model.PasswordHash,
	}, nil
}