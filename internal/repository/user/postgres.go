package user

import (
	"context"
	"errors"

	domain "github.com/21mebrat/lost-found-platform/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

var _ Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Create(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {

	query := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			phone_verified,
			email_verified,
			profile_image_url,
			last_login_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9
		)
		RETURNING
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			phone_verified,
			email_verified,
			profile_image_url,
			status,
			role,
			last_login_at,
			created_at,
			updated_at
	`

	var createdUser domain.User

	err := p.db.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.PhoneVerified,
		user.EmailVerified,
		user.ProfileImageURL,
		user.LastLoginAt,
	).Scan(
		&createdUser.ID,
		&createdUser.FirstName,
		&createdUser.LastName,
		&createdUser.Email,
		&createdUser.Phone,
		&createdUser.PasswordHash,
		&createdUser.PhoneVerified,
		&createdUser.EmailVerified,
		&createdUser.ProfileImageURL,
		&createdUser.Status,
		&createdUser.Role,
		&createdUser.LastLoginAt,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (p *PostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {

	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			phone_verified,
			email_verified,
			profile_image_url,
			status,
			role,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	return p.scanUser(ctx, query, id)
}

func (p *PostgresRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {

	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			phone_verified,
			email_verified,
			profile_image_url,
			status,
			role,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	return p.scanUser(ctx, query, email)
}

func (p *PostgresRepository) GetByPhone(
	ctx context.Context,
	phone string,
) (*domain.User, error) {

	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			phone_verified,
			email_verified,
			profile_image_url,
			status,
			role,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE phone = $1
	`

	return p.scanUser(ctx, query, phone)
}

func (p *PostgresRepository) Update(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {

	query := `
		UPDATE users
		SET
			first_name = $1,
			last_name = $2,
			email = $3,
			phone = $4,
			password_hash = $5,
			phone_verified = $6,
			email_verified = $7,
			profile_image_url = $8,
			status = $9,
			role = $10,
			last_login_at = $11,
			updated_at = $12
		WHERE id = $13
		RETURNING
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			phone_verified,
			email_verified,
			profile_image_url,
			status,
			role,
			last_login_at,
			created_at,
			updated_at
	`

	var updatedUser domain.User

	err := p.db.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.PhoneVerified,
		user.EmailVerified,
		user.ProfileImageURL,
		user.Status,
		user.Role,
		user.LastLoginAt,
		user.UpdatedAt,
		user.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Email,
		&updatedUser.Phone,
		&updatedUser.PasswordHash,
		&updatedUser.PhoneVerified,
		&updatedUser.EmailVerified,
		&updatedUser.ProfileImageURL,
		&updatedUser.Status,
		&updatedUser.Role,
		&updatedUser.LastLoginAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (p *PostgresRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (p *PostgresRepository) scanUser(
	ctx context.Context,
	query string,
	arg any,
) (*domain.User, error) {

	var user domain.User

	err := p.db.QueryRow(
		ctx,
		query,
		arg,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.PhoneVerified,
		&user.EmailVerified,
		&user.ProfileImageURL,
		&user.Status,
		&user.Role,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
