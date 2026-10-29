package registry

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fremenkiel/stdout.cv/internal/session"
)

type sessionCache interface {
	GetFile(id string) (string, error)
}

type Repository struct {
	cache sessionCache
}

func NewRepository(c sessionCache) *Repository {
	return &Repository{
		cache: c,
	}
}

func (r *Repository) GetAll(ctx context.Context, id int, name string) (Registry, error) {
	sessionId, err := session.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	databaseName, err := r.cache.GetFile(sessionId)
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return nil, err
	}

	sql := "SELECT id, current_name FROM table_registry"

	var registry Registry
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table Table
		rows.Scan(&table.Id, &table.CurrentName)

		registry = append(registry, &table)
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return registry, nil
}

func (r *Repository) Create(ctx context.Context, name string) (*int, error) {
	sessionId, err := session.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	databaseName, err := r.cache.GetFile(sessionId)
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("INSERT INTO table_registry (current_name) VALUES ('%s') RETURNING id", name)

	var id int
	if err := db.QueryRowContext(ctx, sql).Scan(&id); err != nil {
		return nil, err
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *Repository) Update(ctx context.Context, id int, name string) error {
	sessionId, err := session.FromContext(ctx)
	if err != nil {
		return err
	}

	databaseName, err := r.cache.GetFile(sessionId)
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("UPDATE table_registry SET current_name = '%s' WHERE id = %v", name, id)

	if _, err := db.ExecContext(ctx, sql); err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Remove(ctx context.Context, id int) error {
	sessionId, err := session.FromContext(ctx)
	if err != nil {
		return err
	}

	databaseName, err := r.cache.GetFile(sessionId)
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("DELETE FROM table_registry WHERE id = %v", id)

	if _, err := db.ExecContext(ctx, sql); err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

