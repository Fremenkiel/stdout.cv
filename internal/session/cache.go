package session

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrFileNotFound = errors.New("session: file not found in cache")

type Cahce struct {
	sessions	sync.Map
	files			map[string]string
}

func NewCache() *Cahce {
	return &Cahce{
		files: make(map[string]string),
	}
}

func (c *Cahce) AddSession(id string) {
	c.sessions.Store(id, time.Now())
	c.AddFile(id)
}

func (c *Cahce) UpdateSession(id string) bool {
	_, loaded := c.sessions.Swap(id, time.Now())
	if !loaded {
		c.AddFile(id)
	}
	return loaded 
}

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
		c.RemoveFile(oldSessions[i])
	}

	return oldSessions
}

func (c *Cahce) AddFile(id string) {
	c.files[id] = fmt.Sprintf("/tmp/stdout_cv_sessions/session_%s.db", id)
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
