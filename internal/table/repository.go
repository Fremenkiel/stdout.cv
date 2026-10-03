package table

import (
	"context"
	"database/sql"

	"github.com/fremenkiel/stdout.cv/internal/platform"
	_ "modernc.org/sqlite"
)

type SessionCache interface {
	GetFile(id string) (string, error)
}

type Repository struct {
	sessionCache	SessionCache
}

func NewRepository(sc SessionCache) *Repository {
	return &Repository{sessionCache: sc}
}

func (r *Repository) GetTables(ctx context.Context) ([]*Table, error) {
	query := `
	SELECT name FROM sqlite_schema WHERE type ='table' AND name NOT LIKE 'sqlite_%';
`

	databaseName, err := r.sessionCache.GetFile(ctx.Value(platform.SessionKey).(string))
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*Table

	for rows.Next() {
		var table Table

		if err := rows.Scan(&table.Name); err != nil {
			return nil, err
		}
		tables = append(tables, &table)
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return tables, nil
}
