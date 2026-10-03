package session

import (
	"testing"
	"time"
)

func TestSessionCache(t *testing.T) {
	tests := []struct {
		name						string
		initialState		map[string]int64
		expectedResult	[]string
	} {
		{
			name: "load_and_delete_two_entries",
			initialState: map[string]int64{
				"1a100bf-c3be-7d9c-acb8-785e8722f892": 0,
				"01a100bf-c3be-7dbb-b038-49c4e6479dba": 4,
				"01a100bf-c3be-7dc3-ac37-5ac1e37cc6b1": 10,
				"01a100bf-c3be-7dc7-abe5-f68536f492a6": 60,
			},
			expectedResult: []string{
				"01a100bf-c3be-7dc7-abe5-f68536f492a6",
				"01a100bf-c3be-7dc3-ac37-5ac1e37cc6b1",
			},
		},
	}

	for i := range tests {
	t.Run(tests[i].name, func(t *testing.T) {
			cache := Cahce{
				files: make(map[string]string),
			}
			
			for si := range tests[i].initialState {
				lastSeen := time.Now().Add(-(time.Duration(tests[i].initialState[si])*time.Minute))
				cache.sessions.Store(si, lastSeen)
				cache.files[si] = "db_name"
			}

			oldSessions := cache.LoadAndDeleteOldSessions()

			if sLen := len(oldSessions); len(tests[i].expectedResult) != sLen {
				t.Fatalf("unexpected return len, got %d", sLen)
			}

			for oi := range oldSessions {
				if !contains(tests[i].expectedResult, oldSessions[oi]) {
					t.Fatalf("deleted session do not match, got %s", oldSessions[oi])
				}
			}

			if fLen := len(cache.files); fLen != len(tests[i].initialState) - len(tests[i].expectedResult) {
				t.Fatalf("unexpected file len after remove, got %d", fLen)
			}
		})
	}
}

func contains(slice []string, value string) bool {
	for i := range slice {
		if slice[i] == value {
			return true
		}
	}
	return false
}
