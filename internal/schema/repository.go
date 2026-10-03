package schema

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *Repository) GetSchema(ctx context.Context, name string) ([]*Schema, error) {
	query := fmt.Sprintf(`
	PRAGMA table_info('%s');
`, name)

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

	var columns []*Schema

	for rows.Next() {
		var column Schema 
		var defaultValue []byte

		if err := rows.Scan(
			&column.Cid,
			&column.Name,
			&column.Type,
			&column.IsNullable,
			&defaultValue,
			&column.IsKey,
			); err != nil {
			return nil, err
		}
		columns = append(columns, &column)
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return columns, nil
}
