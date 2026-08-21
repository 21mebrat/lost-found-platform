package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

	languageCode := user.LanguageCode
	if languageCode == "" {
		languageCode = "en"
	}

	status := user.Status
	if status == "" {
		status = domain.StatusActive
	}

	roleName := user.Role
	if roleName == "" {
		roleName = domain.RoleUser
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO users (
			first_name,
			middle_name,
			last_name,
			phone,
			email,
			fayda,
			language_code,
			password_hash,
			is_phone_verified,
			is_fayda_verified,
			profile_image_url,
			status
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
		RETURNING
			id,
			first_name,
			middle_name,
			last_name,
			phone,
			email,
			fayda,
			language_code,
			password_hash,
			is_phone_verified,
			is_fayda_verified,
			profile_image_url,
			status,
			deleted_at,
			created_at,
			updated_at
	`

	var createdUser domain.User
	err = tx.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.MiddleName,
		user.LastName,
		user.Phone,
		user.Email,
		user.Fayda,
		languageCode,
		user.PasswordHash,
		user.PhoneVerified,
		user.FaydaVerified,
		user.ProfileImageURL,
		status,
	).Scan(
		&createdUser.ID,
		&createdUser.FirstName,
		&createdUser.MiddleName,
		&createdUser.LastName,
		&createdUser.Phone,
		&createdUser.Email,
		&createdUser.Fayda,
		&createdUser.LanguageCode,
		&createdUser.PasswordHash,
		&createdUser.PhoneVerified,
		&createdUser.FaydaVerified,
		&createdUser.ProfileImageURL,
		&createdUser.Status,
		&createdUser.DeletedAt,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Insert role
	roleQuery := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, (SELECT id FROM roles WHERE name = $2))
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err = tx.Exec(ctx, roleQuery, createdUser.ID, string(roleName))
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	createdUser.Role = roleName
	return &createdUser, nil
}

func (p *PostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {

	query := `
		SELECT
			u.id,
			u.first_name,
			u.middle_name,
			u.last_name,
			u.phone,
			u.email,
			u.fayda,
			u.language_code,
			u.password_hash,
			u.is_phone_verified,
			u.is_fayda_verified,
			u.profile_image_url,
			u.status,
			u.deleted_at,
			u.created_at,
			u.updated_at,
			COALESCE(r.name::text, 'user') AS role
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`

	return p.scanUser(ctx, query, id)
}

func (p *PostgresRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {

	query := `
		SELECT
			u.id,
			u.first_name,
			u.middle_name,
			u.last_name,
			u.phone,
			u.email,
			u.fayda,
			u.language_code,
			u.password_hash,
			u.is_phone_verified,
			u.is_fayda_verified,
			u.profile_image_url,
			u.status,
			u.deleted_at,
			u.created_at,
			u.updated_at,
			COALESCE(r.name::text, 'user') AS role
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE LOWER(u.email) = LOWER($1) AND u.deleted_at IS NULL
	`

	return p.scanUser(ctx, query, email)
}

func (p *PostgresRepository) GetByFayda(
	ctx context.Context,
	fayda string,
) (*domain.User, error) {

	query := `
		SELECT
			u.id,
			u.first_name,
			u.middle_name,
			u.last_name,
			u.phone,
			u.email,
			u.fayda,
			u.language_code,
			u.password_hash,
			u.is_phone_verified,
			u.is_fayda_verified,
			u.profile_image_url,
			u.status,
			u.deleted_at,
			u.created_at,
			u.updated_at,
			COALESCE(r.name::text, 'user') AS role
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.fayda = $1 AND u.deleted_at IS NULL
	`

	return p.scanUser(ctx, query, fayda)
}

func (p *PostgresRepository) GetByPhone(
	ctx context.Context,
	phone string,
) (*domain.User, error) {

	query := `
		SELECT
			u.id,
			u.first_name,
			u.middle_name,
			u.last_name,
			u.phone,
			u.email,
			u.fayda,
			u.language_code,
			u.password_hash,
			u.is_phone_verified,
			u.is_fayda_verified,
			u.profile_image_url,
			u.status,
			u.deleted_at,
			u.created_at,
			u.updated_at,
			COALESCE(r.name::text, 'user') AS role
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.phone = $1 AND u.deleted_at IS NULL
	`

	return p.scanUser(ctx, query, phone)
}

func (p *PostgresRepository) Update(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {

	languageCode := user.LanguageCode
	if languageCode == "" {
		languageCode = "en"
	}

	query := `
		WITH updated AS (
			UPDATE users
			SET
				first_name = $1,
				middle_name = $2,
				last_name = $3,
				phone = $4,
				email = $5,
				fayda = $6,
				language_code = $7,
				password_hash = $8,
				is_phone_verified = $9,
				is_fayda_verified = $10,
				profile_image_url = $11,
				status = $12,
				updated_at = NOW()
			WHERE id = $13 AND deleted_at IS NULL
			RETURNING *
		)
		SELECT
			u.id,
			u.first_name,
			u.middle_name,
			u.last_name,
			u.phone,
			u.email,
			u.fayda,
			u.language_code,
			u.password_hash,
			u.is_phone_verified,
			u.is_fayda_verified,
			u.profile_image_url,
			u.status,
			u.deleted_at,
			u.created_at,
			u.updated_at,
			COALESCE(r.name::text, 'user') AS role
		FROM updated u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
	`

	userRes, err := p.scanUser(
		ctx,
		query,
		user.FirstName,
		user.MiddleName,
		user.LastName,
		user.Phone,
		user.Email,
		user.Fayda,
		languageCode,
		user.PasswordHash,
		user.PhoneVerified,
		user.FaydaVerified,
		user.ProfileImageURL,
		user.Status,
		user.ID,
	)

	if err != nil {
		return nil, err
	}
	if userRes == nil {
		return nil, domain.ErrUserNotFound
	}

	return userRes, nil
}

func (p *PostgresRepository) GetList(
	ctx context.Context,
	params ListParams,
) (*PaginatedUsers, error) {

	page := params.Page
	if page < 1 {
		page = 1
	}

	limit := params.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	var sb strings.Builder
	sb.WriteString(`
		SELECT
			u.id,
			u.first_name,
			u.middle_name,
			u.last_name,
			u.phone,
			u.email,
			u.fayda,
			u.language_code,
			u.password_hash,
			u.is_phone_verified,
			u.is_fayda_verified,
			u.profile_image_url,
			u.status,
			u.deleted_at,
			u.created_at,
			u.updated_at,
			COALESCE(r.name::text, 'user') AS role,
			COUNT(*) OVER() AS total_count
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.deleted_at IS NULL
	`)

	args := make([]any, 0, 4)
	argIdx := 1

	if search := strings.TrimSpace(params.Search); search != "" {
		pattern := "%" + search + "%"
		sb.WriteString(fmt.Sprintf(`
			AND (
				first_name ILIKE $%d OR
				middle_name ILIKE $%d OR
				last_name ILIKE $%d OR
				phone ILIKE $%d OR
				email ILIKE $%d OR
				fayda ILIKE $%d
			)
		`, argIdx, argIdx, argIdx, argIdx, argIdx, argIdx))
		args = append(args, pattern)
		argIdx++
	}

	if params.Status != nil {
		sb.WriteString(fmt.Sprintf(" AND status = $%d", argIdx))
		args = append(args, *params.Status)
		argIdx++
	}

	sortBy := "created_at"
	switch strings.ToLower(params.SortBy) {
	case "first_name", "last_name", "email", "phone", "status", "created_at":
		sortBy = strings.ToLower(params.SortBy)
	}

	sortOrder := "DESC"
	if strings.ToUpper(params.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	sb.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder))
	sb.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1))
	args = append(args, limit, offset)

	rows, err := p.db.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	var totalCount int64

	for rows.Next() {
		var u domain.User
		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.MiddleName,
			&u.LastName,
			&u.Phone,
			&u.Email,
			&u.Fayda,
			&u.LanguageCode,
			&u.PasswordHash,
			&u.PhoneVerified,
			&u.FaydaVerified,
			&u.ProfileImageURL,
			&u.Status,
			&u.DeletedAt,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.Role,
			&totalCount,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	totalPages := 0
	if limit > 0 && totalCount > 0 {
		totalPages = int((totalCount + int64(limit) - 1) / int64(limit))
	}

	return &PaginatedUsers{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}

func (p *PostgresRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	query := `
		UPDATE users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
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
	args ...any,
) (*domain.User, error) {

	var user domain.User

	err := p.db.QueryRow(
		ctx,
		query,
		args...,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.Fayda,
		&user.LanguageCode,
		&user.PasswordHash,
		&user.PhoneVerified,
		&user.FaydaVerified,
		&user.ProfileImageURL,
		&user.Status,
		&user.DeletedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
