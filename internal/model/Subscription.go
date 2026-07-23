package model

import "time"

// Subscription представляет подписку пользователя на монету.
type Subscription struct {
	ID              int64
	UserID          int64
	CoinID          int64
	IntervalMinutes int
	LastSentAt      time.Time
}

// NewSubscription создает новую подписку.
func NewSubscription(userID, coinID int64) *Subscription {
	return &Subscription{
		UserID:          userID,
		CoinID:          coinID,
		IntervalMinutes: 10,
		LastSentAt:      time.Now(),
	}
}

// NewSubscriptionFromDB восстанавливает подписку из БД.
func NewSubscriptionFromDB(
	id int64,
	userID int64,
	coinID int64,
	intervalMinutes int,
	lastSentAt time.Time,
) *Subscription {
	return &Subscription{
		ID:              id,
		UserID:          userID,
		CoinID:          coinID,
		IntervalMinutes: intervalMinutes,
		LastSentAt:      lastSentAt,
	}
}

// SetInterval устанавливает интервал в минутах.
func (s *Subscription) SetInterval(minutes int) {
	if minutes < 1 {
		minutes = 1
	}
	s.IntervalMinutes = minutes
}

// IsReadyToSend проверяет, пора ли отправлять уведомление.
func (s *Subscription) IsReadyToSend() bool {
	return time.Since(s.LastSentAt) >= s.Interval()
}

// Interval возвращает интервал как time.Duration.
func (s *Subscription) Interval() time.Duration {
	return time.Duration(s.IntervalMinutes) * time.Minute
}

// MarkSent обновляет время последней отправки.
func (s *Subscription) MarkSent() {
	s.LastSentAt = time.Now()
}
