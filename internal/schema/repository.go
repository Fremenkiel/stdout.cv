package schema

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/fremenkiel/stdout.cv/internal/platform"
	"github.com/fremenkiel/stdout.cv/internal/table"
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

func (r *Repository) GetSchemas(ctx context.Context, tables []*table.Table) ([]*table.Table, error) {
	fkQuerySlice := make([]string, len(tables))
	querySlice := make([]string, len(tables))

	for i := range tables {
		fkQuerySlice[i] = fmt.Sprintf(`
			SELECT "from" FROM pragma_foreign_key_list('%s');
			`, tables[i].Name)
		querySlice[i] = fmt.Sprintf(`
			SELECT * FROM pragma_table_info('%s');
			`, tables[i].Name)
	}

	fkQuery := strings.Join(fkQuerySlice, "")
	query := strings.Join(querySlice, "")

	databaseName, err := r.sessionCache.GetFile(ctx.Value(platform.SessionKey).(string))
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

	for i := range tables {
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

		for rows.Next() {
			var column table.Column
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

			log.Printf("%s: table %s", column.Name, tables[i].Name)

			_, ok := fkMap[column.Name]
			column.IsForeignKey = ok

			tables[i].Columns = append(tables[i].Columns, &column)
		}
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *Repository) GetSchema(ctx context.Context, name string) ([]*table.Column, error) {
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

	var columns []*table.Column

	for rows.Next() {
		var column table.Column 
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
		columns = append(columns, &column)
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return columns, nil
}
