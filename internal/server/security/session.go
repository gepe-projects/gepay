package security

// ! simpen dulu aja skrg fokus mvp dulu
// package session

// import (
// 	"context"
// 	"crypto/rand"
// 	"encoding/base64"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// var (
// 	ErrSessionNotFound    = errors.New("session not found")
// 	ErrSessionExpired     = errors.New("session expired")
// 	ErrInvalidSession     = errors.New("invalid session")
// 	ErrSessionRevoked     = errors.New("session revoked from another device")
// 	ErrGenerateSessionID  = errors.New("failed to generate session id")
// )

// // Session represents user session data
// type Session struct {
// 	ID        string    `json:"id"`
// 	UserID    string    `json:"user_id"`
// 	Role      string    `json:"role"`
// 	CreatedAt time.Time `json:"created_at"`
// 	ExpiresAt time.Time `json:"expires_at"`
// 	IPAddress string    `json:"ip_address,omitempty"`
// 	UserAgent string    `json:"user_agent,omitempty"`
// }

// // Manager handles session operations
// type Manager struct {
// 	redis      *redis.Client
// 	prefix     string
// 	ttl        time.Duration
// 	idleTimeout time.Duration
// }

// // Config holds session manager configuration
// type Config struct {
// 	Prefix      string        // Redis key prefix (default: "session")
// 	TTL         time.Duration // Session time to live (default: 24 hours)
// 	IdleTimeout time.Duration // Idle timeout for session refresh (default: 1 hour)
// }

// // NewManager creates a new session manager with best practice defaults
// func NewManager(redisClient *redis.Client, config Config) *Manager {
// 	// Set defaults
// 	if config.Prefix == "" {
// 		config.Prefix = "session"
// 	}
// 	if config.TTL == 0 {
// 		config.TTL = 24 * time.Hour
// 	}
// 	if config.IdleTimeout == 0 {
// 		config.IdleTimeout = 1 * time.Hour
// 	}

// 	return &Manager{
// 		redis:       redisClient,
// 		prefix:      config.Prefix,
// 		ttl:         config.TTL,
// 		idleTimeout: config.IdleTimeout,
// 	}
// }

// // generateSecureID generates cryptographically secure session ID
// func (m *Manager) generateSecureID() (string, error) {
// 	// 32 bytes = 256 bits of entropy
// 	b := make([]byte, 32)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return "", fmt.Errorf("%w: %v", ErrGenerateSessionID, err)
// 	}
// 	// Base64 URL encoding untuk cookie-safe string
// 	return base64.URLEncoding.EncodeToString(b), nil
// }

// // sessionKey returns Redis key for session data
// func (m *Manager) sessionKey(sessionID string) string {
// 	return fmt.Sprintf("%s:data:%s", m.prefix, sessionID)
// }

// // userSessionKey returns Redis key for user-to-session mapping
// func (m *Manager) userSessionKey(userID string) string {
// 	return fmt.Sprintf("%s:user:%s", m.prefix, userID)
// }

// // Create creates a new session and revokes any existing session for the user
// // This enforces single session per user policy
// func (m *Manager) Create(ctx context.Context, userID, role, ipAddress, userAgent string) (string, error) {
// 	// Step 1: Revoke existing session (single session enforcement)
// 	if err := m.revokeUserSession(ctx, userID); err != nil {
// 		// Log warning but don't fail - old session might not exist
// 		// In production, use proper logger here
// 	}

// 	// Step 2: Generate secure session ID
// 	sessionID, err := m.generateSecureID()
// 	if err != nil {
// 		return "", err
// 	}

// 	// Step 3: Create session object
// 	now := time.Now()
// 	sess := Session{
// 		ID:        sessionID,
// 		UserID:    userID,
// 		Role:      role,
// 		CreatedAt: now,
// 		ExpiresAt: now.Add(m.ttl),
// 		IPAddress: ipAddress,
// 		UserAgent: userAgent,
// 	}

// 	// Step 4: Serialize session data
// 	data, err := json.Marshal(sess)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal session: %w", err)
// 	}

// 	// Step 5: Store session in Redis with TTL
// 	sessionKey := m.sessionKey(sessionID)
// 	if err := m.redis.Set(ctx, sessionKey, data, m.ttl).Err(); err != nil {
// 		return "", fmt.Errorf("failed to store session: %w", err)
// 	}

// 	// Step 6: Store user-to-session mapping
// 	userKey := m.userSessionKey(userID)
// 	if err := m.redis.Set(ctx, userKey, sessionID, m.ttl).Err(); err != nil {
// 		// Rollback: delete session data if mapping fails
// 		m.redis.Del(ctx, sessionKey)
// 		return "", fmt.Errorf("failed to store user mapping: %w", err)
// 	}

// 	return sessionID, nil
// }

// // Get retrieves session data and validates it
// func (m *Manager) Get(ctx context.Context, sessionID string) (*Session, error) {
// 	if sessionID == "" {
// 		return nil, ErrInvalidSession
// 	}

// 	// Retrieve session data
// 	sessionKey := m.sessionKey(sessionID)
// 	data, err := m.redis.Get(ctx, sessionKey).Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			return nil, ErrSessionNotFound
// 		}
// 		return nil, fmt.Errorf("failed to get session: %w", err)
// 	}

// 	// Deserialize session
// 	var sess Session
// 	if err := json.Unmarshal([]byte(data), &sess); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
// 	}

// 	// Check if session is expired
// 	if time.Now().After(sess.ExpiresAt) {
// 		m.Delete(ctx, sessionID)
// 		return nil, ErrSessionExpired
// 	}

// 	// Validate session ownership (anti-session hijacking)
// 	userKey := m.userSessionKey(sess.UserID)
// 	currentSessionID, err := m.redis.Get(ctx, userKey).Result()
// 	if err != nil || currentSessionID != sessionID {
// 		// User logged in from another device
// 		return nil, ErrSessionRevoked
// 	}

// 	return &sess, nil
// }

// // Refresh extends session lifetime (for idle timeout management)
// func (m *Manager) Refresh(ctx context.Context, sessionID string) error {
// 	// Get current session
// 	sess, err := m.Get(ctx, sessionID)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if session needs refresh (idle timeout logic)
// 	timeUntilExpiry := time.Until(sess.ExpiresAt)
// 	if timeUntilExpiry > m.idleTimeout {
// 		// Session still has plenty of time, no need to refresh
// 		return nil
// 	}

// 	// Update expiration time
// 	sess.ExpiresAt = time.Now().Add(m.ttl)

// 	// Serialize updated session
// 	data, err := json.Marshal(sess)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal session: %w", err)
// 	}

// 	// Update session in Redis with new TTL
// 	sessionKey := m.sessionKey(sessionID)
// 	if err := m.redis.Set(ctx, sessionKey, data, m.ttl).Err(); err != nil {
// 		return fmt.Errorf("failed to refresh session: %w", err)
// 	}

// 	// Update user mapping TTL
// 	userKey := m.userSessionKey(sess.UserID)
// 	if err := m.redis.Expire(ctx, userKey, m.ttl).Err(); err != nil {
// 		return fmt.Errorf("failed to refresh user mapping: %w", err)
// 	}

// 	return nil
// }

// // Delete removes session and user mapping
// func (m *Manager) Delete(ctx context.Context, sessionID string) error {
// 	// Get session to find userID
// 	sess, err := m.Get(ctx, sessionID)
// 	if err != nil {
// 		// Session might not exist, but we should still try to delete
// 		if err == ErrSessionNotFound {
// 			return nil
// 		}
// 		// For other errors, still attempt cleanup
// 	}

// 	// Delete session data
// 	sessionKey := m.sessionKey(sessionID)
// 	m.redis.Del(ctx, sessionKey)

// 	// Delete user mapping if we have userID
// 	if sess != nil {
// 		userKey := m.userSessionKey(sess.UserID)
// 		m.redis.Del(ctx, userKey)
// 	}

// 	return nil
// }

// // DeleteByUserID deletes all sessions for a user (logout from all devices)
// func (m *Manager) DeleteByUserID(ctx context.Context, userID string) error {
// 	return m.revokeUserSession(ctx, userID)
// }

// // revokeUserSession revokes existing session for a user
// func (m *Manager) revokeUserSession(ctx context.Context, userID string) error {
// 	userKey := m.userSessionKey(userID)

// 	// Get existing session ID
// 	oldSessionID, err := m.redis.Get(ctx, userKey).Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			// No existing session, nothing to revoke
// 			return nil
// 		}
// 		return fmt.Errorf("failed to get user session: %w", err)
// 	}

// 	// Delete old session data
// 	oldSessionKey := m.sessionKey(oldSessionID)
// 	m.redis.Del(ctx, oldSessionKey)

// 	// Delete user mapping
// 	m.redis.Del(ctx, userKey)

// 	return nil
// }

// // Validate checks if session exists and is valid (lightweight check)
// func (m *Manager) Validate(ctx context.Context, sessionID string) (bool, error) {
// 	_, err := m.Get(ctx, sessionID)
// 	if err != nil {
// 		if errors.Is(err, ErrSessionNotFound) ||
// 		   errors.Is(err, ErrSessionExpired) ||
// 		   errors.Is(err, ErrSessionRevoked) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return true, nil
// }

// // GetUserID returns userID from session without full validation
// // Useful for logging/analytics where you don't need full session validation
// func (m *Manager) GetUserID(ctx context.Context, sessionID string) (string, error) {
// 	sess, err := m.Get(ctx, sessionID)
// 	if err != nil {
// 		return "", err
// 	}
// 	return sess.UserID, nil
// }

// // CountActiveSessions returns number of active sessions (for monitoring)
// func (m *Manager) CountActiveSessions(ctx context.Context) (int64, error) {
// 	pattern := fmt.Sprintf("%s:data:*", m.prefix)
// 	keys, err := m.redis.Keys(ctx, pattern).Result()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int64(len(keys)), nil
// }
