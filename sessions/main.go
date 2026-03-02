package main

import (
	"fmt"
	"time"
)

type Session struct {
	ID        string
	CreatedAt time.Time
}

type SessionManager struct {
	sessions map[string]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (m *SessionManager) Create(id string) *Session {
	s := &Session{
		ID:        id,
		CreatedAt: time.Now(),
	}
	m.sessions[id] = s
	return s
}

func (m *SessionManager) Get(id string) *Session {
	return m.sessions[id]
}

func (m *SessionManager) Cleanup(ttl time.Duration) {
	for id, s := range m.sessions {
		if time.Since(s.CreatedAt) > ttl {
			delete(m.sessions, id)
		}
	}
}

func main() {
	m := NewSessionManager()

	go func() {
		for {
			m.Cleanup(5 * time.Second)
			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < 1000; i++ {
		go func(i int) {
			m.Create(fmt.Sprintf("%d", i))
		}(i)
	}

	time.Sleep(10 * time.Second)
}
