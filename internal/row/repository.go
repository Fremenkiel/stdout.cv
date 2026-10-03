package row

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
	sessionCache SessionCache
}

func NewRepository(sc SessionCache) *Repository {
	return &Repository{sessionCache: sc}
}

func (r *Repository) GetRows(ctx context.Context, query string, fields []string, table string) ([]*Row, error) {
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

	var list []*Row

	for rows.Next() {
		rawBuffer := make([]sql.RawBytes, len(fields))
		callArgs := make([]interface{}, len(rawBuffer))

		for i := range rawBuffer {
			callArgs[i] = &rawBuffer[i]
		}

		if err := rows.Scan(callArgs...); err != nil {
			return nil, err
		}

		var row Row
		for i := range rawBuffer {
			str := string(rawBuffer[i])
			row = append(row, &str)
		}

		list = append(list, &row)
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	return list, nil
}
