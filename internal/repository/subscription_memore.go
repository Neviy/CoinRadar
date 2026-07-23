// package repository хранит реализацию репозитория для работы с подписками.
package repository

import (
	"CoinRadar/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SubscriptionRepository реализует репозиторий для работы с подписками.
type SubscriptionRepository struct {
	conn *pgxpool.Pool
}

// NewSubscriptionRepository создает новый экземпляр SubscriptionRepository.
func NewSubscriptionRepository(conn *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		conn: conn,
	}
}

// Create сохраняет новую подписку.
func (sr *SubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	const query = `
		INSERT INTO subscriptions (
			user_id,
			coin_id,
			interval_minutes,
			last_sent_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var id int64
	err := sr.conn.QueryRow(
		ctx,
		query,
		subscription.UserID,
		subscription.CoinID,
		subscription.IntervalMinutes,
		subscription.LastSentAt,
	).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	subscription.ID = id
	return nil
}

// GetByUserID возвращает все подписки пользователя.
func (sr *SubscriptionRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Subscription, error) {
	const query = `
		SELECT
			id,
			user_id,
			coin_id,
			interval_minutes,
			last_sent_at
		FROM subscriptions
		WHERE user_id = $1;
	`
	rows, err := sr.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()
	var subscriptions []*model.Subscription
	for rows.Next() {
		var (
			id              int64
			userIDDB        int64
			coinID          int64
			intervalMinutes int
			lastSentAt      time.Time
		)
		if err := rows.Scan(
			&id,
			&userIDDB,
			&coinID,
			&intervalMinutes,
			&lastSentAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, model.NewSubscriptionFromDB(
			id,
			userIDDB,
			coinID,
			intervalMinutes,
			lastSentAt,
		))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// GetByUserAndCoin возвращает подписку пользователя на конкретную монету.
func (sr *SubscriptionRepository) GetByUserAndCoin(ctx context.Context, userID, coinID int64) (*model.Subscription, error) {
	const query = `
		SELECT
			id,
			user_id,
			coin_id,
			interval_minutes,
			last_sent_at
		FROM subscriptions
		WHERE user_id = $1
		AND coin_id = $2;
	`
	var (
		id              int64
		userIDDB        int64
		coinIDDB        int64
		intervalMinutes int
		lastSentAt      time.Time
	)
	err := sr.conn.QueryRow(ctx, query, userID, coinID).Scan(
		&id,
		&userIDDB,
		&coinIDDB,
		&intervalMinutes,
		&lastSentAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return model.NewSubscriptionFromDB(
		id,
		userIDDB,
		coinIDDB,
		intervalMinutes,
		lastSentAt,
	), nil
}

// Update обновляет подписку.
func (sr *SubscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	if subscription.ID == 0 {
		return errors.New("subscription ID is required")
	}
	const query = `
		UPDATE subscriptions
		SET
			interval_minutes = $1,
			last_sent_at = $2
		WHERE id = $3;
	`
	result, err := sr.conn.Exec(
		ctx,
		query,
		subscription.IntervalMinutes,
		subscription.LastSentAt,
		subscription.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

// Delete удаляет подписку пользователя на монету.
func (sr *SubscriptionRepository) Delete(ctx context.Context, userID, coinID int64) error {
	const query = `
		DELETE FROM subscriptions
		WHERE user_id = $1
		AND coin_id = $2;
	`
	result, err := sr.conn.Exec(ctx, query, userID, coinID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("subscription not found")
	}
	return nil
}
