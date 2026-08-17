package token

import (
	"context"
	"errors"

	tokendomain "github.com/21mebrat/lost-found-platform/internal/domain/token"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

var _ Repository = (*PostgresRepository)(nil)

// constructor
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) Create(ctx context.Context, token *tokendomain.RefreshToken) error {
	query := `
			INSERT INTO refresh_tokens(
			id,
			user_id,
			token,
			revoked,
			expiresAt
			createdAt
			)
			VALUES(
			$1,$2,$3,FALSE,$5,NOW()
			)
			`

	_, err := p.db.Exec(
		ctx,
		query,
		uuid.New(),
		token.UserID,
		token.Token,
		token.ExpiresAt,
	)

	return err

}

func (p *PostgresRepository) Get(ctx context.Context, value string) (*tokendomain.RefreshToken, error) {
	query := `
   SELECT 
    id ,
	user_id,
	token,
	revoked,
	created_at,
	expired_at
	FROM refresh_tokens 
	WHERE token= $1
`

	var token tokendomain.RefreshToken

	err := p.db.QueryRow(
		ctx,
		query,
		value,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.CreatedAt,
		&token.ExpiresAt,
		&token.Revoked,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}
func (p *PostgresRepository) Revoke(ctx context.Context, value string) error {
	query := `

	UPDATE refresh_tokens

	SET revoked=true

	WHERE token=$1

	`
	_, err := p.db.Exec(
		ctx,
		query,
		value,
	)

	return err
}
func (p *PostgresRepository) Delete(ctx context.Context, value string) error {
	query := `
	DELETE FROM refresh_tokens
	WHERE token=$1
	`

	_, err := p.db.Exec(
		ctx,
		query,
		value,
	)

	return err
}

func (p *PostgresRepository) DeleteByUserId(ctx context.Context, user_id uuid.UUID) error {
	query := `
	DELETE FROM refresh_tokens
	WHERE token=$1
	`

	_, err := p.db.Exec(
		ctx,
		query,
		user_id,
	)

	return err
}
