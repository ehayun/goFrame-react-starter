package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tzlev/internal/redis"
)

type Session struct {
	UserID    uint      `json:"user_id"` // Legacy field, not used
	Zehut     string    `json:"zehut"`   // Israeli ID - primary key in users table
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionManager struct {
	prefix string
	ttl    time.Duration
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		prefix: "tzlev:session:",
		ttl:    24 * time.Hour,
	}
}

func (sm *SessionManager) key(sessionID string) string {
	return fmt.Sprintf("%s%s", sm.prefix, sessionID)
}

func (sm *SessionManager) Create(ctx context.Context, sessionID string, session *Session) error {
	session.CreatedAt = time.Now()

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := sm.key(sessionID)
	return redis.Client.Set(ctx, key, data, sm.ttl).Err()
}

func (sm *SessionManager) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := sm.key(sessionID)

	data, err := redis.Client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (sm *SessionManager) Delete(ctx context.Context, sessionID string) error {
	key := sm.key(sessionID)
	return redis.Client.Del(ctx, key).Err()
}

func (sm *SessionManager) Refresh(ctx context.Context, sessionID string) error {
	key := sm.key(sessionID)
	return redis.Client.Expire(ctx, key, sm.ttl).Err()
}
