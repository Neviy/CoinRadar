package model

import (
	"errors"
	"strings"
	"time"
)

// Coin представляет криптовалюту.
type Coin struct {
	ID          int64     `json:"id"`
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Recommended bool      `json:"recommended"`
	Price       float64   `json:"price"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewCoin конструктор для инициализации.
func NewCoin(name string, price float64, symbol string) (*Coin, error) {
	var errs []error

	if name == "" {
		errs = append(errs, errors.New("invalid name"))
	}

	if price < 0 {
		errs = append(errs, errors.New("price cannot be less than 0"))
	}

	symbol = strings.ToUpper(symbol)
	if symbol == "" {
		errs = append(errs, errors.New("invalid symbol"))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Coin{
		Name:        name,
		Price:       price,
		Recommended: false,
		Symbol:      symbol,
		UpdatedAt:   time.Now(),
	}, nil
}

// NewCoinFromDB восстановление монеты из БД.
func NewCoinFromDB(
	id int64,
	name string,
	symbol string,
	recommended bool,
	price float64,
	updatedAt time.Time,
) *Coin {
	return &Coin{
		ID:          id,
		Name:        name,
		Symbol:      symbol,
		Recommended: recommended,
		Price:       price,
		UpdatedAt:   updatedAt,
	}
}

// UpdatePrice обновляет цену монеты.
func (c *Coin) UpdatePrice(price float64) error {
	if price < 0 {
		return errors.New("price cannot be less than 0")
	}
	c.Price = price
	c.UpdatedAt = time.Now()
	return nil
}

// MinutesSinceUpdate возвращает количество минут, прошедших с момента последнего обновления  монеты.
func (c *Coin) MinutesSinceUpdate() int {
	if c.UpdatedAt.IsZero() {
		return 0
	}
	diff := time.Since(c.UpdatedAt)
	if diff < 0 {
		return 0
	}
	return int(diff / time.Minute)
}

// MarkRecommended помечает монету как рекомендованную.
func (c *Coin) MarkRecommended(){
	c.Recommended = true
}

// UnmarkRecommended снимает пометку с монеты как рекомендованной.
func (c *Coin) UnmarkRecommended(){
	c.Recommended = false
}
