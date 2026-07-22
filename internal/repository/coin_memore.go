// package repository хранит реализацию репозиториев для работы с монетами.
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

//CoinRepository реализует репозиторий для работы с монетами.
type CoinRepository struct{
	conn *pgxpool.Pool
}

//NewCoinRepository создает новый экземпляр CoinRepository 
func NewCoinRepository(conn *pgxpool.Pool)*CoinRepository{
	return &CoinRepository{
		conn: conn,
	}
}

// SaveCoin сохраняет монету в базе данных.
func (cr *CoinRepository) SaveCoin(ctx context.Context, coin *model.Coin) error {
	const query = `
		INSERT INTO coins (
			symbol,
			name,
			price,
			recommended,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var id int64
	err := cr.conn.QueryRow(
		ctx,
		query,
		coin.Symbol,
		coin.Name,
		coin.Price,
		coin.Recommended,
		coin.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}
	coin.ID = id
	return nil
}

// GetCoinByID возвращает монету по её ID.
func (cr *CoinRepository) GetCoinByID(ctx context.Context, coinID int64) (*model.Coin, error) {
	const query = `
		SELECT
			id,
			symbol,
			name,
			recommended,
			price,
			updated_at
		FROM coins
		WHERE id = $1;
	`
	var (
		id          int64
		symbol      string
		name        string
		recommended bool
		price       float64
		updatedAt   time.Time
	)
	err := cr.conn.QueryRow(ctx, query, coinID).Scan(
		&id,
		&symbol,
		&name,
		&recommended,
		&price,
		&updatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get coin: %w", err)
	}
	return model.NewCoinFromDB(
		id,
		name,
		symbol,
		recommended,
		price,
		updatedAt,
	), nil
}

// GetCoinBySymbol возвращает монету по символу.
func (cr *CoinRepository) GetCoinBySymbol(ctx context.Context, symbol string) (*model.Coin, error) {
	const query = `
		SELECT
			id,
			symbol,
			name,
			recommended,
			price,
			updated_at
		FROM coins
		WHERE symbol = $1;
	`
	var (
		id          int64
		symbolDB    string
		name        string
		recommended bool
		price       float64
		updatedAt   time.Time
	)
	err := cr.conn.QueryRow(ctx, query, symbol).Scan(
		&id,
		&symbolDB,
		&name,
		&recommended,
		&price,
		&updatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get coin: %w", err)
	}
	return model.NewCoinFromDB(
		id,
		name,
		symbolDB,
		recommended,
		price,
		updatedAt,
	), nil
}

// GetAllCoins возвращает все монеты.
func (cr *CoinRepository) GetAllCoins(ctx context.Context) ([]*model.Coin, error) {
	const query = `
		SELECT
			id,
			symbol,
			name,
			recommended,
			price,
			updated_at
		FROM coins;
	`
	rows, err := cr.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get coins: %w", err)
	}
	defer rows.Close()
	var coins []*model.Coin
	for rows.Next() {
		var (
			id          int64
			symbol      string
			name        string
			recommended bool
			price       float64
			updatedAt   time.Time
		)
		if err := rows.Scan(
			&id,
			&symbol,
			&name,
			&recommended,
			&price,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan coin: %w", err)
		}
		coins = append(coins, model.NewCoinFromDB(
			id,
			name,
			symbol,
			recommended,
			price,
			updatedAt,
		))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return coins, nil
}

// UpdateCoin обновляет данные монеты.
func (cr *CoinRepository) UpdateCoin(ctx context.Context, coin *model.Coin) error {
	if coin.ID == 0 {
		return errors.New("coin ID is required")
	}
	const query = `
		UPDATE coins
		SET
			price = $1,
			recommended = $2,
			updated_at = $3
		WHERE id = $4;
	`
	result, err := cr.conn.Exec(
		ctx,
		query,
		coin.Price,
		coin.Recommended,
		coin.UpdatedAt,
		coin.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update coin: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("coin not found")
	}
	return nil
}

// DeleteCoin удаляет монету по символу.
func (cr *CoinRepository) DeleteCoin(ctx context.Context, symbol string) error {
	if symbol == "" {
		return errors.New("coin symbol is required")
	}
	const query = `
		DELETE FROM coins
		WHERE symbol = $1;
	`
	result, err := cr.conn.Exec(ctx, query, symbol)
	if err != nil {
		return fmt.Errorf("failed to delete coin: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("coin not found")
	}
	return nil
}