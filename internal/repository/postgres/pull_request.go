package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/rand"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/errors/repository"
	"time"
)

const MAX_REVIEWERS = 2

type pullRequestRepositoryPostgres struct {
	pool     *pgxpool.Pool
	teamRepo *teamRepositoryPostgres
	userRepo *userRepositoryPostgres
}

func NewPullRequestRepositoryPostgres(pool *pgxpool.Pool, teamRepo *teamRepositoryPostgres, userRepo *userRepositoryPostgres) *pullRequestRepositoryPostgres {
	return &pullRequestRepositoryPostgres{
		pool:     pool,
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (r *pullRequestRepositoryPostgres) CreateWithReviewers(ctx context.Context, pr domain.PullRequest) (*domain.PullRequestWithReviewers, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	author, err := r.userRepo.GetByIDTx(ctx, tx, pr.AuthorID)
	if err != nil {
		// ошибка не оборачивается, так как ошибка оборачивается на предыдущем уровне
		return nil, err
	}

	teamWithMembers, err := r.teamRepo.GetTeamWithMembersTx(ctx, tx, author.TeamName)
	if err != nil {
		return nil, err
	}

	// Случайный выбор ревьюеров
	activeMembers := make([]string, 0, len(teamWithMembers.Members))
	for _, u := range teamWithMembers.Members {
		if u.ID != pr.AuthorID && u.IsActive {
			activeMembers = append(activeMembers, u.ID)
		}
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(activeMembers), func(i, j int) {
		activeMembers[i], activeMembers[j] = activeMembers[j], activeMembers[i]
	})

	if len(activeMembers) > MAX_REVIEWERS {
		activeMembers = activeMembers[:MAX_REVIEWERS]
	}

	queryCreatePR := `
        INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
        VALUES ($1, $2, $3, 'OPEN')
    `

	_, err = tx.Exec(ctx, queryCreatePR, pr.ID, pr.Name, pr.AuthorID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, repository.ErrPullRequestExists
			}
		}
		return nil, repository.ErrInternalError
	}

	querySetReviewers := `
            INSERT INTO pr_reviewers (pull_request_id, user_id)
            VALUES ($1, $2)
        `

	for _, reviewerID := range activeMembers {
		_, err := tx.Exec(ctx, querySetReviewers, pr.ID, reviewerID)
		if err != nil {
			return nil, repository.ErrInternalError
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, repository.ErrInternalError
	}

	return &domain.PullRequestWithReviewers{
		PullRequest:       pr,
		AssignedReviewers: activeMembers,
	}, nil
}

func (r *pullRequestRepositoryPostgres) Merge(ctx context.Context, prID string) (*domain.PullRequestWithReviewers, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	existing, err := r.getPRWithReviewersTx(ctx, tx, prID)
	if err != nil {
		return nil, err
	}

	// проверка для идемпотентности
	if existing.Status == domain.PRStatusMerged {
		if err = tx.Commit(ctx); err != nil {
			return nil, repository.ErrInternalError
		}
		return existing, nil
	}

	queryMerge := `
        UPDATE pull_requests
        SET status = 'MERGED', merged_at = NOW()
        WHERE pull_request_id = $1
    `
	_, err = tx.Exec(ctx, queryMerge, prID)
	if err != nil {
		return nil, repository.ErrInternalError
	}

	updated, err := r.getPRWithReviewersTx(ctx, tx, prID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, repository.ErrInternalError
	}

	return updated, nil
}

func (r *pullRequestRepositoryPostgres) getPRWithReviewersTx(ctx context.Context, tx pgx.Tx, prID string) (*domain.PullRequestWithReviewers, error) {
	queryGetPR := `
        SELECT pull_request_id, pull_request_name, author_id, status, merged_at
        FROM pull_requests
        WHERE pull_request_id = $1
    `
	var pr domain.PullRequest

	err := tx.QueryRow(ctx, queryGetPR, prID).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.MergedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, repository.ErrPullRequestNotFound
		default:
			return nil, repository.ErrInternalError
		}
	}

	queryReviewers := `
        SELECT user_id
        FROM pr_reviewers
        WHERE pull_request_id = $1
    `

	rows, err := tx.Query(ctx, queryReviewers, prID)
	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer rows.Close()

	reviewers := []string{}
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, repository.ErrInternalError
		}
		reviewers = append(reviewers, uid)
	}

	return &domain.PullRequestWithReviewers{
		PullRequest:       pr,
		AssignedReviewers: reviewers,
	}, nil
}

func (r *pullRequestRepositoryPostgres) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.Reviewer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	pr, err := r.getPRWithReviewersTx(ctx, tx, prID)
	if err != nil {
		return nil, err
	}

	// проверка на MERGED
	if pr.Status == domain.PRStatusMerged {
		return nil, repository.ErrPullRequestMerged
	}

	// проверка есть ли вообще такой пользователь
	_, err = r.userRepo.GetByIDTx(ctx, tx, oldReviewerID)
	if err != nil {
		return nil, err
	}

	// проверка, что oldReviewer действительно был назначен
	found := false
	for _, uid := range pr.AssignedReviewers {
		if uid == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, repository.ErrReviewerNotAssigned
	}

	oldUser, err := r.userRepo.GetByIDTx(ctx, tx, oldReviewerID)
	if err != nil {
		return nil, err
	}

	team, err := r.teamRepo.GetTeamWithMembersTx(ctx, tx, oldUser.TeamName)
	if err != nil {
		return nil, err
	}

	candidates := []string{}
	for _, u := range team.Members {
		if u.ID != oldReviewerID && u.ID != pr.AuthorID && u.IsActive {
			candidates = append(candidates, u.ID)
		}
	}

	// проверка на наличие кандидата
	if len(candidates) == 0 {
		return nil, repository.ErrNoReplacementCandidate
	}

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	newReviewer := candidates[randGen.Intn(len(candidates))]

	queryUpdate := `
        UPDATE pr_reviewers
        SET user_id = $1
        WHERE pull_request_id = $2 AND user_id = $3
    `
	_, err = tx.Exec(ctx, queryUpdate, newReviewer, prID, oldReviewerID)
	if err != nil {
		return nil, repository.ErrInternalError
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, repository.ErrInternalError
	}

	return &domain.Reviewer{
		ID: newReviewer,
	}, nil
}
