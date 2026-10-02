package user

import (
	"context"
	"database/sql"
)

type Repository struct {
	db	*sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUser(ctx context.Context) (*User, error) {
	query := `
	SELECT name FROM users ORDER BY id LIMIT 1;
`

	user := &User{}

	err := r.db.QueryRowContext(ctx, query).Scan(
		&user.Name,
		)

	if err != nil {
		return nil, err
	}

	return user, nil
}
