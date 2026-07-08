// package services хранит бизнес-логику связанную с монетами.
package services

import (
	"CoinRadar/internal/model"
	"context"
	"errors"
	"strings"
)

// CoinRepository определяет интерфейс для работы с монетами.
type CoinRepository interface {
	GetCoinBySymbol(ctx context.Context, symbol string) (*model.Coin, error)
	SaveCoin(ctx context.Context, coin *model.Coin) error
	GetAllCoins(ctx context.Context) ([]*model.Coin, error)
	UpdateCoin(ctx context.Context, coin *model.Coin) error
	DeleteCoin(ctx context.Context, symbol string) error
}


// CoinService предоставляет методы для работы с монетами.
type CoinService struct {
	coinRepo CoinRepository
}

// NewCoinService создает новый экземпляр CoinService.
func NewCoinService(coinRepo CoinRepository) *CoinService {
	return &CoinService{
		coinRepo: coinRepo,
	}
}

//CreateCoin создает новую монету.
func (s *CoinService) CreateCoin(ctx context.Context, name string, price float64, symbol string) (*model.Coin, error) {
	symbol = strings.ToUpper(symbol)
	if coin, _ := s.coinRepo.GetCoinBySymbol(ctx, symbol); coin != nil {
		return nil, errors.New("coin with this symbol already exists")
	}
	coin, err := model.NewCoin(name, price, symbol)
	if err != nil {
		return nil, err
	}
	err = s.coinRepo.SaveCoin(ctx, coin)
	if err != nil {
		return nil, err
	}
	return coin, nil
}

// GetCoinBySymbol возвращает монету по символу.
func (s *CoinService) GetCoinBySymbol(ctx context.Context, symbol string) (*model.Coin, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	symbol = strings.ToUpper(symbol)
	coin, err := s.coinRepo.GetCoinBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return coin, nil
}

// GetAllCoins возвращает все монеты.
func (s *CoinService) GetAllCoins(ctx context.Context) ([]*model.Coin, error) {
	coins, err := s.coinRepo.GetAllCoins(ctx)
	if err != nil {
		return nil, err
	}
	return coins, nil
}

// UpdateCoin обновляет информацию о монете.
func (s *CoinService) UpdateCoin(ctx context.Context, coin *model.Coin) error {
	if coin == nil {
		return errors.New("coin cannot be nil")
	}
	err := s.coinRepo.UpdateCoin(ctx, coin)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCoin удаляет монету по символу.
func (s *CoinService) DeleteCoin(ctx context.Context, symbol string) error {
	if symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	symbol = strings.ToUpper(symbol)
	err := s.coinRepo.DeleteCoin(ctx, symbol)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCoinPrice обновляет цену монеты по символу.
func (s *CoinService) UpdateCoinPrice(ctx context.Context, symbol string, newPrice float64) error {
	if symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	symbol = strings.ToUpper(symbol)
	coin, err := s.coinRepo.GetCoinBySymbol(ctx, symbol)
	if err != nil {
		return err
	}
	if err := coin.UpdatePrice(newPrice); err != nil {
		return err
	}
	err = s.coinRepo.UpdateCoin(ctx, coin)
	if err != nil {
		return err
	}
	return nil
}

// GetRecommendedCoins возвращает список рекомендованных монет.
func (s *CoinService) GetRecommendedCoins(ctx context.Context) ([]*model.Coin, error) {
	coins, err := s.coinRepo.GetAllCoins(ctx)
	if err != nil {
		return nil, err
	}
	var recommendedCoins []*model.Coin
	for _, coin := range coins {
		if coin.IsRecommended() {
			recommendedCoins = append(recommendedCoins, coin)
		}
	}
	return recommendedCoins, nil
}

// MarkCoinAsRecommended помечает монету как рекомендованную.
func (s *CoinService) MarkCoinAsRecommended(ctx context.Context, symbol string) error {
	if symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	symbol = strings.ToUpper(symbol)
	coin, err := s.coinRepo.GetCoinBySymbol(ctx, symbol)
	if err != nil {
		return err
	}
	coin.MarkRecommended()
	return s.coinRepo.UpdateCoin(ctx, coin)
}

// UnmarkCoinAsRecommended снимает пометку с монеты как рекомендованной.
func (s *CoinService) UnmarkCoinAsRecommended(ctx context.Context, symbol string) error {
	if symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	symbol = strings.ToUpper(symbol)
	coin, err := s.coinRepo.GetCoinBySymbol(ctx, symbol)
	if err != nil {
		return err
	}
	coin.UnmarkRecommended()
	return s.coinRepo.UpdateCoin(ctx, coin)
}

//GetOutdatedCoins возвращает список устаревших монет, которые не обновлялись.
func (s *CoinService) GetOutdatedCoins(ctx context.Context, minutes int) ([]*model.Coin, error) {
	coins, err := s.coinRepo.GetAllCoins(ctx)
	if err != nil {
		return nil, err
	}
	var outdatedCoins []*model.Coin
	for _, coin := range coins {
		if coin.MinutesSinceUpdate() > minutes {
			outdatedCoins = append(outdatedCoins, coin)
		}
	}
	return outdatedCoins, nil
}

