package schema

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/session"
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

func (r *Repository) GetSchema(ctx context.Context, name string) ([]*column.Column, error) {
	fkQuery := fmt.Sprintf(`
			SELECT "from" FROM pragma_foreign_key_list('%s');
			`, name)
	query := fmt.Sprintf(`
			SELECT * FROM pragma_table_info('%s');
			`, name)

	sessionId, err := session.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	databaseName, err := r.sessionCache.GetFile(sessionId)
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_foreign_keys=on", databaseName))
	if err != nil {
		return nil, err
	}

	fkRows, err := db.QueryContext(ctx, fkQuery)
	if err != nil {
		return nil, err
	}
	defer fkRows.Close()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fkMap := make(map[string]struct{})
	for fkRows.Next() {
		var columnName string
		if err := fkRows.Scan(
			&columnName,
			); err != nil {
			return nil, err
		}

		fkMap[columnName] = struct{}{}
	}

	var columns []*column.Column
	for rows.Next() {
		var column column.Column
		var defaultValue []byte

		if err := rows.Scan(
			&column.Cid,
			&column.Name,
			&column.Type,
			&column.IsNullable,
			&defaultValue,
			&column.IsPrimaryKey,
			); err != nil {
			return nil, err
		}

		_, ok := fkMap[column.Name]
		column.IsForeignKey = ok

		columns = append(columns, &column)
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return columns, nil
}

