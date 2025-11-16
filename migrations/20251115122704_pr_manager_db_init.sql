-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams (
                       team_name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE users (
                       user_id VARCHAR(255) PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       team_name VARCHAR(255) REFERENCES teams(team_name) ON DELETE SET NULL,
                       is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE pull_requests (
                               pull_request_id VARCHAR(255) PRIMARY KEY,
                               pull_request_name VARCHAR(255) NOT NULL,
                               author_id VARCHAR(255) REFERENCES users(user_id) ON DELETE CASCADE,
                               status pr_status NOT NULL DEFAULT 'OPEN',
                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               merged_at TIMESTAMPTZ
);

CREATE TABLE pr_reviewers (
                              pull_request_id VARCHAR(255) REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
                              user_id VARCHAR(255) REFERENCES users(user_id) ON DELETE CASCADE,
                              PRIMARY KEY (pull_request_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pr_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TYPE IF EXISTS pr_status;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
