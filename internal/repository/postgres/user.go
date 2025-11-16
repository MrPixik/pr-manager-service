package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/errors/repository"
)

type userRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewUserRepositoryPostgres(pool *pgxpool.Pool) *userRepositoryPostgres {
	return &userRepositoryPostgres{pool: pool}
}

func (r *userRepositoryPostgres) UpsertManyTx(ctx context.Context, tx pgx.Tx, users []domain.User) error {
	query := `
        INSERT INTO users (user_id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE
        SET username = EXCLUDED.username,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active
    `
	for _, u := range users {
		if _, err := tx.Exec(ctx, query, u.ID, u.Username, u.TeamName, u.IsActive); err != nil {
			return repository.ErrInternalError
		}
	}
	return nil
}

func (r *userRepositoryPostgres) GetByTeamName(ctx context.Context, teamName string) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, username, team_name, is_active FROM users WHERE team_name=$1`,
		teamName,
	)
	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, repository.ErrInternalError
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *userRepositoryPostgres) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	query := `
        UPDATE users
        SET is_active = $1
        WHERE user_id = $2
        RETURNING user_id, username, team_name, is_active
    `

	var u domain.User
	err := r.pool.QueryRow(ctx, query, isActive, userID).Scan(
		&u.ID, &u.Username, &u.TeamName, &u.IsActive,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repository.ErrUserNotFound
		default:
			return nil, repository.ErrInternalError
		}
	}

	return &u, nil
}

func (r *userRepositoryPostgres) GetByIDTx(ctx context.Context, tx pgx.Tx, userID string) (*domain.User, error) {
	query := `
        SELECT user_id, username, team_name, is_active
        FROM users
        WHERE user_id = $1
    `

	var u domain.User
	err := tx.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Username, &u.TeamName, &u.IsActive,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repository.ErrUserNotFound
		default:
			return nil, repository.ErrInternalError
		}
	}

	return &u, nil
}

func (r *userRepositoryPostgres) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
        SELECT user_id, username, team_name, is_active
        FROM users
        WHERE user_id = $1
    `

	var u domain.User
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Username, &u.TeamName, &u.IsActive,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repository.ErrUserNotFound
		default:
			return nil, repository.ErrInternalError
		}
	}

	return &u, nil
}

func (r *userRepositoryPostgres) GetReviewPullRequests(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	// проверка на существование такого пользователя
	if _, err := r.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	queryGetReviewPR := `
        SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
        FROM pull_requests pr
        JOIN pr_reviewers r ON r.pull_request_id = pr.pull_request_id
        WHERE r.user_id = $1
        ORDER BY pr.created_at DESC
    `

	rows, err := r.pool.Query(ctx, queryGetReviewPR, userID)
	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer rows.Close()

	var prs []domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return nil, repository.ErrInternalError
		}
		prs = append(prs, pr)
	}

	return prs, nil
}
