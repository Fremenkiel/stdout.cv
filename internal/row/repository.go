package row

import (
	"context"
	"database/sql"

	"github.com/fremenkiel/stdout.cv/internal/query"
	"github.com/fremenkiel/stdout.cv/internal/session"
	_ "modernc.org/sqlite"
)

type SessionCache interface {
	GetFile(id string) (string, error)
}

type Repository struct {
	sessionCache SessionCache
}

var _ repository = (*Repository)(nil)

func NewRepository(sc SessionCache) *Repository {
	return &Repository{sessionCache: sc}
}

func (r *Repository) GetRows(ctx context.Context, queryString string, query *query.Query) ([]*Row, error) {
	sessionId, err := session.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	databaseName, err := r.sessionCache.GetFile(sessionId)
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Row

	for rows.Next() {
		rawBuffer := make([]sql.RawBytes, len(query.Fields))
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
