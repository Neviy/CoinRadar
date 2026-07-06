package model

import "time"

// Subscription хранит информацию о подписке.
type Subscription struct {
	ID         int64         `json:"id"`
	UserID     int64         `json:"user_id"`
	CoinID     int64         `json:"coin_id"`
	Interval   time.Duration `json:"interval"`
	LastSentAt time.Time     `json:"last_sent_at"`
}

// NewSubscription создает новую подписку.
func NewSubscription(userID, coinID int64) *Subscription {
	return &Subscription{
		UserID:     userID,
		CoinID:     coinID,
		Interval:   10 * time.Minute,
		LastSentAt: time.Now(),
	}
}

// NewSubscriptionFromDB создает подписку из БД.
func NewSubscriptionFromDB(
	id, userID, coinID int64,
	interval time.Duration,
	lastSentAt time.Time,
) *Subscription {
	return &Subscription{
		ID:         id,
		UserID:     userID,
		CoinID:     coinID,
		Interval:   interval,
		LastSentAt: lastSentAt,
	}
}

// SetInterval устанавливает интервал.
func (s *Subscription) SetInterval(interval time.Duration) {
	if interval < time.Minute {
		interval = time.Minute
	}
	s.Interval = interval
}

// IsReadyToSend проверяет, пора ли отправлять уведомление.
func (s *Subscription) IsReadyToSend() bool {
	return time.Since(s.LastSentAt) >= s.Interval
}

// MarkSent обновляет время последней отправки.
func (s *Subscription) MarkSent() {
	s.LastSentAt = time.Now()
}