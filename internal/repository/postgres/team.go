package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/errors/repository"
)

type teamRepositoryPostgres struct {
	userRepo *userRepositoryPostgres
	pool     *pgxpool.Pool
}

func NewTeamRepositoryPostgres(pool *pgxpool.Pool, userRepo *userRepositoryPostgres) *teamRepositoryPostgres {
	return &teamRepositoryPostgres{
		userRepo: userRepo,
		pool:     pool,
	}
}

func (r *teamRepositoryPostgres) AddTeamWithMembers(ctx context.Context, team domain.Team, members []domain.User) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	if err = r.InsertTx(ctx, tx, team); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err = r.userRepo.UpsertManyTx(ctx, tx, members); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *teamRepositoryPostgres) GetTeamWithMembers(ctx context.Context, teamName string) (*domain.TeamWithUsers, error) {
	team, err := r.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	members, err := r.userRepo.GetByTeamName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	return &domain.TeamWithUsers{
		TeamName: team.Name,
		Members:  members,
	}, nil
}

// ручка проверяет существование команды
func (r *teamRepositoryPostgres) GetTeamWithMembersTx(ctx context.Context, tx pgx.Tx, teamName string) (*domain.TeamWithUsers, error) {
	var team domain.Team
	err := tx.QueryRow(ctx, `SELECT team_name FROM teams WHERE team_name=$1`, teamName).Scan(&team.Name)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repository.ErrTeamNotFound
		default:
			return nil, repository.ErrInternalError
		}
	}

	rows, err := tx.Query(ctx,
		`SELECT user_id, username, team_name, is_active FROM users WHERE team_name=$1`,
		teamName,
	)
	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer rows.Close()

	var members []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, repository.ErrInternalError
		}
		members = append(members, u)
	}

	return &domain.TeamWithUsers{
		TeamName: team.Name,
		Members:  members,
	}, nil
}

func (r *teamRepositoryPostgres) InsertTx(ctx context.Context, tx pgx.Tx, team domain.Team) error {
	sql := `
        INSERT INTO teams (team_name)
        VALUES ($1)
    `

	_, err := tx.Exec(ctx, sql, team.Name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return repository.ErrTeamAlreadyExists
			}
			return repository.ErrInternalError
		}
	}
	return err
}

func (r *teamRepositoryPostgres) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	queryGetTeam := `SELECT team_name
		FROM teams
		WHERE team_name=$1
		`
	var team domain.Team
	err := r.pool.QueryRow(ctx, queryGetTeam, name).
		Scan(&team.Name)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repository.ErrTeamNotFound
		default:
			return nil, repository.ErrInternalError
		}
	}
	return &team, nil
}

// GetTeamStats возвращает статистику по команде:
// количество активных и неактивных пользователей,
// количество открытых и замерженных PR
func (r *teamRepositoryPostgres) GetTeamStats(ctx context.Context, teamName string) (activeUsers, inactiveUsers, openPRs, mergedPRs int, err error) {
	sql := `
	SELECT 
		COUNT(DISTINCT u.user_id) FILTER (WHERE u.is_active) AS active_users,
		COUNT(DISTINCT u.user_id) FILTER (WHERE NOT u.is_active) AS inactive_users,
		COUNT(pr.pull_request_id) FILTER (WHERE pr.status='OPEN') AS open_prs,
		COUNT(pr.pull_request_id) FILTER (WHERE pr.status='MERGED') AS merged_prs
	FROM teams t
	LEFT JOIN users u ON u.team_name = t.team_name
	LEFT JOIN pull_requests pr ON pr.author_id = u.user_id
	WHERE t.team_name = $1
	GROUP BY t.team_name
	`

	row := r.pool.QueryRow(ctx, sql, teamName)
	err = row.Scan(&activeUsers, &inactiveUsers, &openPRs, &mergedPRs)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return 0, 0, 0, 0, repository.ErrTeamNotFound
		default:
			return 0, 0, 0, 0, repository.ErrInternalError
		}
	}

	return activeUsers, inactiveUsers, openPRs, mergedPRs, nil
}
