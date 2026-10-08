package session

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/platform/database"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
	"github.com/google/uuid"
)

func TestCreateSession(t *testing.T) {
	cache := NewCache(database.TestDatabaseFileNameTemplate)
	service := NewService(cache)

	id, err := service.CreateSession()
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(t.Context(), id)
	defer testutil.CleanupSession(ctx)

	if _, ok := cache.sessions.Load(id); !ok {
		t.Fatal("session not saved in cache")
	}

	if _, ok := cache.files[id]; !ok {
		t.Fatal("session database file path not saved in cache")
	}
}

func TestServiceUpdateSession(t *testing.T) {
	tests := []struct{
		name			string
		loaded		bool
	} {
		{
			name: "update_existing",
			loaded: true,
		},
		{
			name: "create_missing",
			loaded: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := NewCache(database.TestDatabaseFileNameTemplate)
			service := NewService(cache)

			sessionUuid, err := uuid.NewV7()
			if err != nil {
				t.Fatal(err)
			}

			id := sessionUuid.String()

			if test.loaded {
				cache.sessions.Store(id, nil)
				cache.files[id] = "placeholder_value"
			}

			if err := service.UpdateSession(id); err != nil {
				t.Fatal(err)
			}
			ctx := NewContext(t.Context(), id)
			defer testutil.CleanupSession(ctx)

			if _, ok := cache.sessions.Load(id); !ok {
				t.Fatal("session not saved in cache")
			}

			if _, ok := cache.files[id]; !ok {
				t.Fatal("session database file path not saved in cache")
			}
		})
	}
}

func TestRemoveExpiredSessions(t *testing.T) {
	tests := []struct{
		name					string
		initialState			map[string]int16	
		sessionLoaded	bool
		fileLoaded		bool
		fileCreated		bool
		expectedError	error
	} {
		{
			name: "remove_session",
			initialState: map[string]int16{
				"1a100bf-c3be-7d9c-acb8-785e8722f892": 0,
				"01a100bf-c3be-7dbb-b038-49c4e6479dba": 4,
				"01a100bf-c3be-7dc3-ac37-5ac1e37cc6b1": 2,
				"01a100bf-c3be-7dc7-abe5-f68536f492a6": 1,
			},
			sessionLoaded: true,
			fileLoaded: true,
			fileCreated: true,
			expectedError: nil,
		},
		{
			name: "session_not_loaded",
			initialState: map[string]int16{
				"1a100bf-c3be-7d9c-acb8-785e8722f892": 0,
				"01a100bf-c3be-7dbb-b038-49c4e6479dba": 4,
				"01a100bf-c3be-7dc3-ac37-5ac1e37cc6b1": 2,
				"01a100bf-c3be-7dc7-abe5-f68536f492a6": 1,
			},
			sessionLoaded: false,
			fileLoaded: true,
			fileCreated: true,
			expectedError: nil,
		},
		{
			name: "file_not_created",
			initialState: map[string]int16{
				"1a100bf-c3be-7d9c-acb8-785e8722f892": 0,
				"01a100bf-c3be-7dbb-b038-49c4e6479dba": 4,
				"01a100bf-c3be-7dc3-ac37-5ac1e37cc6b1": 2,
				"01a100bf-c3be-7dc7-abe5-f68536f492a6": 1,
			},
			sessionLoaded: true,
			fileLoaded: true,
			fileCreated: false,
			expectedError: nil,
		},
		{
			name: "file_not_loaded",
			initialState: map[string]int16{
				"1a100bf-c3be-7d9c-acb8-785e8722f892": 0,
				"01a100bf-c3be-7dbb-b038-49c4e6479dba": 4,
				"01a100bf-c3be-7dc3-ac37-5ac1e37cc6b1": 2,
				"01a100bf-c3be-7dc7-abe5-f68536f492a6": 1,
			},
			sessionLoaded: true,
			fileLoaded: false,
			fileCreated: true,
			expectedError: nil,
		},
		{
			name: "file_not_loaded_remove_all",
			initialState: map[string]int16{
				"1a100bf-c3be-7d9c-acb8-785e8722f892": 0,
				"01a100bf-c3be-7dbb-b038-49c4e6479dba": 4,
				"01a100bf-c3be-7dc3-ac37-5ac1e37cc6b1": 2,
				"01a100bf-c3be-7dc7-abe5-f68536f492a6": 10,
			},
			sessionLoaded: true,
			fileLoaded: false,
			fileCreated: true,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := NewCache(database.TestDatabaseFileNameTemplate)
			service := NewService(cache)

			var expired []string
			
			for key, value := range test.initialState {
				lastSeen := time.Now().Add(-(time.Duration(value)*time.Minute))
				cache.sessions.Store(key, lastSeen)
				if value > 5 {
					expired = append(expired, key)
				}
			}

			sessionUuid, err := uuid.NewV7()
			if err != nil {
				t.Fatal(err)
			}

			id := sessionUuid.String()
			fileName := fmt.Sprintf(database.TestDatabaseFileNameTemplate, id)

			expired = append(expired, id)

			if test.sessionLoaded {
				cache.sessions.Store(id, time.Now().Add(-(6*time.Minute)))
			}

			if test.fileLoaded {
				cache.files[id] = fileName
			}

			if test.fileCreated {
				if err := createTestFile(fileName); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := ensureFolder(); err != nil {
					t.Fatal(err)
				}
			}

			if err := service.RemoveExpiredSessions(); !errors.Is(err, test.expectedError) {
				t.Fatalf("unexpected error, %v got %v", test.expectedError, err)
			}

			if _, err := os.Stat(fileName); err == nil {
				t.Fatal("file was not removed")
			}

			for _, session := range expired {
				if _, ok := cache.sessions.Load(session); ok {
					t.Fatalf("session not removed from cache, id %s", session)
				}

				if _, ok := cache.files[session]; ok {
					t.Fatalf("session file path not removed from cache, id %s", session)
				}
			}
		})
	}
}

func TestCreateSessionDatabase(t *testing.T) {
	cache := NewCache(database.TestDatabaseFileNameTemplate)
	service := NewService(cache)

	if err := ensureFolder(); err != nil {
		t.Fatal(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}

	if err := service.createSessionDatabase(id.String()); err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(t.Context(), id.String())
	defer testutil.CleanupSession(ctx)

	if len(cache.files) != 1 {
		t.Fatalf("unexpected file slice len, 1 got %d", len(cache.files))
	}

	databaseName, err := cache.GetFile(id.String())
	if err != nil {
		t.Fatalf("unexpected error while getting file, %v", err)
	}

	if _, err := os.Stat(databaseName); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveSessionDatabase(t *testing.T) {
	cache := NewCache(database.TestDatabaseFileNameTemplate)
	service := NewService(cache)

	id, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}

	databaseName := cache.AddFile(id.String())

	if err := createTestFile(databaseName); err != nil {
		t.Fatal(err)
	}

	if err := service.removeSessionDatabase(id.String()); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(databaseName); err == nil {
		t.Fatal("file was not removed")
	}

	if len(cache.files) != 0 {
		t.Fatal("cache file slice is not empty")
	}
}

func createTestFile(name string) error {
	if err := ensureFolder(); err != nil {
		return err
	}

	dstFileHandler, err := os.OpenFile(name, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer dstFileHandler.Close()

	if _, err = dstFileHandler.Write([]byte{1}); err != nil {
		return err
	}

	return nil
}

func ensureFolder() error {
	_, err := os.Stat("/tmp/stdout_cv_sessions")
	if err != nil {
		if err := os.MkdirAll("/tmp/stdout_cv_sessions", os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

