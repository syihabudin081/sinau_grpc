package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"user-service/pb"
)

type AuthProvider struct {
	db *gorm.DB
}

func NewAuthProvider(db *gorm.DB) *AuthProvider {
	return &AuthProvider{db: db}
}

func (p *AuthProvider) CreateUser(ctx context.Context, user *pb.UserORM) (string, error) {
	// Check if username or email already exists
	var existingUser pb.UserORM
	result := p.db.WithContext(ctx).Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser)
	if result.Error == nil {
		// User already exists
		if existingUser.Username == user.Username {
			return "", errors.New("username already exists")
		}
		return "", errors.New("email already exists")
	}

	// Create new user
	if err := p.db.WithContext(ctx).Create(user).Error; err != nil {
		return "", err
	}

	return strconv.FormatUint(user.Id, 10), nil
}

func (p *AuthProvider) FindUserByUsername(ctx context.Context, username string) (*pb.UserORM, error) {
	var user pb.UserORM
	if err := p.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (p *AuthProvider) FindUserByID(ctx context.Context, userID string) (*pb.UserORM, error) {
	var user pb.UserORM
	if err := p.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (p *AuthProvider) UpdateUserStatus(ctx context.Context, userID uint64, isActive bool) error {
	result := p.db.WithContext(ctx).Model(&pb.UserORM{}).Where("id = ?", userID).Update("is_active", isActive)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
