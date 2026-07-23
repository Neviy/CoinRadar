// SubscriptionServices package services хранит бизнес-логику, связанную с подписками.
package services

import (
	"CoinRadar/internal/model"
	"context"
	"errors"
	"strings"
)

// SubscriptionRepository определяет интерфейс для работы с подписками.
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *model.Subscription) error
	GetByUserID(ctx context.Context, userID int64) ([]*model.Subscription, error)
	GetByUserAndCoin(ctx context.Context, userID, coinID int64) (*model.Subscription, error)
	Update(ctx context.Context, subscription *model.Subscription) error
	Delete(ctx context.Context, userID, coinID int64) error
}

// SubscriptionService предоставляет методы для работы с подписками.
type SubscriptionService struct {
	subscriptionRepo SubscriptionRepository
	coinService      *CoinService
	userService      *UserService
}

// NewSubscriptionService создает новый экземпляр SubscriptionService.
func NewSubscriptionService(subscriptionRepo SubscriptionRepository, coinService *CoinService, userService *UserService) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		coinService:      coinService,
		userService:      userService,
	}
}

// SubscribeUserToCoin подписывает пользователя на монету.
func (s *SubscriptionService) SubscribeUserToCoin(ctx context.Context, userID int64, coinSymbol string) error {
	symbol := strings.ToUpper(coinSymbol)
	if user, err := s.userService.GetUserByID(ctx, userID); err != nil || user == nil {
		return errors.New("user not found")
	}
	coin, err := s.coinService.GetCoinBySymbol(ctx, symbol)
	if err != nil || coin == nil {
		return errors.New("coin not found")
	}
	if subscription, _ := s.subscriptionRepo.GetByUserAndCoin(ctx, userID, coin.ID); subscription != nil {
		return errors.New("user is already subscribed to this coin")
	}
	subscription := model.NewSubscription(userID, coin.ID)
	err = s.subscriptionRepo.Create(ctx, subscription)
	if err != nil {
		return err
	}
	return nil
}

// UnsubscribeUserFromCoin отписывает пользователя от монеты.
func (s *SubscriptionService) UnsubscribeUserFromCoin(ctx context.Context, userID int64, coinSymbol string) error {
	symbol := strings.ToUpper(coinSymbol)
	if user, err := s.userService.GetUserByID(ctx, userID); err != nil || user == nil {
		return errors.New("user not found")
	}
	coin, err := s.coinService.GetCoinBySymbol(ctx, symbol)
	if err != nil || coin == nil {
		return errors.New("coin not found")
	}
	subscription, err := s.subscriptionRepo.GetByUserAndCoin(ctx, userID, coin.ID)
	if err != nil || subscription == nil {
		return errors.New("subscription not found")
	}
	err = s.subscriptionRepo.Delete(ctx, userID, coin.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetUserSubscriptions возвращает список подписок пользователя.
func (s *SubscriptionService) GetUserSubscriptions(ctx context.Context, userID int64) ([]*model.Subscription, error) {
	if user, err := s.userService.GetUserByID(ctx, userID); err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	subscriptions, err := s.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// IsUserSubscribedToCoin проверяет, подписан ли пользователь на монету.
func (s *SubscriptionService) IsUserSubscribedToCoin(ctx context.Context, userID int64, coinSymbol string) (bool, error) {
	symbol := strings.ToUpper(coinSymbol)
	if user, err := s.userService.GetUserByID(ctx, userID); err != nil || user == nil {
		return false, errors.New("user not found")
	}
	coin, err := s.coinService.GetCoinBySymbol(ctx, symbol)
	if err != nil || coin == nil {
		return false, errors.New("coin not found")
	}
	subscription, err := s.subscriptionRepo.GetByUserAndCoin(ctx, userID, coin.ID)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, nil
	}
	return true, nil
}

// UpdateSubscription обновляет данные подписки.
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, subscription *model.Subscription) error {
	if subscription == nil {
		return errors.New("subscription cannot be nil")
	}
	if user, err := s.userService.GetUserByID(ctx, subscription.UserID); err != nil || user == nil {
		return errors.New("user not found")
	}
	if coin, err := s.coinService.GetCoinByID(ctx, subscription.CoinID); err != nil || coin == nil {
		return errors.New("coin not found")
	}
	err := s.subscriptionRepo.Update(ctx, subscription)
	if err != nil {
		return err
	}
	return nil
}
