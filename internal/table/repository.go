package table

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

func (r *Repository) GetTables(ctx context.Context) ([]*Table, error) {
	query := `
	SELECT name FROM pragma_table_list WHERE type = 'table' AND name NOT LIKE 'sqlite_%';
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

		countQuery := fmt.Sprintf(`
			SELECT COUNT(*) FROM %s;
			`, table.Name)

		if err := db.QueryRowContext(ctx, countQuery).Scan(&table.RowCount); err != nil {
			return nil, err
		}

		tables = append(tables, &table)
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return tables, nil
}
