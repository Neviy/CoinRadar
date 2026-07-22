// package repository хранит реализацию репозиториев для работы с юзерами.
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

// UserRepository реализует репозиторий для работы с пользователями.
type UserRepository struct {
	conn *pgxpool.Pool
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

// SaveUser сохраняет пользователя.
func (r *UserRepository) SaveUser(ctx context.Context, user *model.User) error {
	const query = `
		INSERT INTO users (
			telegram_id,
			email,
			birthday,
			premium,
			admin,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id;
	`
	err := r.conn.QueryRow(
		ctx,
		query,
		user.TelegramID,
		user.Email,
		user.Birthday,
		user.Premium,
		user.Admin,
		user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// GetUserByTelegramID возвращает пользователя по Telegram ID.
func (r *UserRepository) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	const query = `
		SELECT
			id,
			telegram_id,
			email,
			birthday,
			premium,
			admin,
			created_at
		FROM users
		WHERE telegram_id = $1;
	`
	var (
		id        int64
		email     *string
		birthday  time.Time
		premium   bool
		admin     bool
		createdAt time.Time
	)
	err := r.conn.QueryRow(ctx, query, telegramID).Scan(
		&id,
		&telegramID,
		&email,
		&birthday,
		&premium,
		&admin,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return model.NewUserFromDB(
		id,
		telegramID,
		email,
		birthday,
		premium,
		admin,
		createdAt,
	), nil
}

// GetUserByID возвращает пользователя по ID.
func (r *UserRepository) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	const query = `
		SELECT
			id,
			telegram_id,
			email,
			birthday,
			premium,
			admin,
			created_at
		FROM users
		WHERE id = $1;
	`
	var (
		telegramID int64
		email      *string
		birthday   time.Time
		premium    bool
		admin      bool
		createdAt  time.Time
	)
	err := r.conn.QueryRow(ctx, query, userID).Scan(
		&userID,
		&telegramID,
		&email,
		&birthday,
		&premium,
		&admin,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return model.NewUserFromDB(
		userID,
		telegramID,
		email,
		birthday,
		premium,
		admin,
		createdAt,
	), nil
}

// UpdateUser обновляет пользователя.
func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	const query = `
		UPDATE users
		SET
			telegram_id = $1,
			email = $2,
			birthday = $3,
			premium = $4,
			admin = $5
		WHERE id = $6;
	`
	result, err := r.conn.Exec(
		ctx,
		query,
		user.TelegramID,
		user.Email,
		user.Birthday,
		user.Premium,
		user.Admin,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// DeleteUser удаляет пользователя.
func (r *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	const query = `
		DELETE FROM users
		WHERE id = $1;
	`
	result, err := r.conn.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}