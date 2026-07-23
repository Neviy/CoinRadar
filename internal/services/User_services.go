// package services хранит бизнес-логику, связанную с пользователями.
package services

import (
	"CoinRadar/internal/model"
	"context"
	"errors"
	"time"
)

// UserRepository определяет интерфейс для работы с пользователями.
type UserRepository interface {
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	SaveUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, userID int64) error
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
}

// UserService предоставляет методы для работы с пользователями.
type UserService struct {
	userRepo UserRepository
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser создает нового пользователя.
func (s *UserService) CreateUser(ctx context.Context, telegramID int64, birthday time.Time) error {
	if user, err := s.userRepo.GetUserByTelegramID(ctx, telegramID); err != nil || user != nil {
		return errors.New("user with this telegram ID already exists")
	}
	user, err := model.NewUser(telegramID, birthday)
	if err != nil {
		return err
	}
	if err := s.userRepo.SaveUser(ctx, user); err != nil {
		return err
	}
	return nil
}

// GetUserByID возвращает пользователя по его ID.
func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	if userID <= 0 {
		return nil, errors.New("user id must be positive")
	}
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser обновляет данные пользователя.
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	if user.ID <= 0 {
		return errors.New("user id must be positive")
	}
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

// DeleteUser удаляет пользователя по его ID.
func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errors.New("user id must be positive")
	}
	return s.userRepo.DeleteUser(ctx, userID)
}

// SetUserEmail устанавливает адрес электронной почты пользователя.
func (s *UserService) SetUserEmail(ctx context.Context, userID int64, email string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	user.SetEmail(email)
	return s.userRepo.UpdateUser(ctx, user)
}

//CheckStatusUser проверяет статус пользователя.
func (s *UserService) CheckStatusUser(ctx context.Context, userID int64) (bool, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("user not found")
	}
	return user.Premium, nil
}

//MakeAdmin - делает пользователя администратором.
func (s *UserService) MakeAdmin(ctx context.Context, telegramID int64) (*model.User, error) {
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.PromoteToAdmin()
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// RemoveAdmin - удаляет права администратора у пользователя.
func (s *UserService) RemoveAdmin(ctx context.Context, telegramID int64) (*model.User, error) {
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if !user.Admin {
		return nil, errors.New("user is not an admin")
	}
	user.PromoteToUser()
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

//MakePremium - делает пользователя премиум.
func (s *UserService) MakePremium(ctx context.Context, telegramID int64) (*model.User, error) {
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.MakePremium()
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

//RemovePremium - удаляет премиум статус у пользователя.
func (s *UserService) RemovePremium(ctx context.Context, telegramID int64) (*model.User, error) {
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if !user.Premium {
		return nil, errors.New("user is not premium")
	}
	user.RemovePremium()
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
