package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/platform/database"
)

type Cahce struct {
	sessions	sync.Map
	files			map[string]string
}

func NewCache() *Cahce {
	return &Cahce{
		files: make(map[string]string),
	}
}

// Adds session id to cache, alongside timestamp.
// Uses Update under the hood, and ingores response
func (c *Cahce) AddSession(id string) {
	c.UpdateSession(id)
}

// Loads session by id. If no session is found, then an empty struct is returned
func (c *Cahce) LoadSession(id string) (time.Time, bool) {
	if value, ok := c.sessions.Load(id); ok {
		return value.(time.Time), ok
	}
	return time.Time{}, false
}

// Updated session timestamp in cache, adds if not exists.
// Returns whether or not the session already were present in the cache.
func (c *Cahce) UpdateSession(id string) bool {
	_, loaded := c.sessions.Swap(id, time.Now())

	return loaded 
}

// Deletes all sessions older then 5 minutes from the cache.
// Returns the deleted session ids.
func (c *Cahce) LoadAndDeleteOldSessions() []string {
	var oldSessions []string

	c.sessions.Range(func(key any, value any) bool {
		if value.(time.Time).Before(time.Now().Add(-(5*time.Minute))) {
			oldSessions = append(oldSessions, key.(string))
		}
		return true
	})

	for i := range oldSessions {
		c.sessions.Delete(oldSessions[i])
	}

	return oldSessions
}

// Adds session database file path to cache
func (c *Cahce) AddFile(id string) string {
	filePath := fmt.Sprintf(database.DatabaseFileNameTemplate, id)

	c.files[id] = filePath

	return filePath
}

func (c *Cahce) RemoveFile(id string) {
	delete(c.files, id)
}

func (c *Cahce) GetFile(id string) (string, error) {
	str, ok := c.files[id]
	if !ok {
		return "", ErrFileNotFound
	}

	return str, nil
}
