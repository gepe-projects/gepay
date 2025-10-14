package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	authAttemptsKey = "auth_attempts"
	authBanKey      = "auth_ban"
)

// CheckBan: apakah user sedang diban?
func (s *Security) CheckBan(ctx context.Context, id string) (time.Duration, error) {
	banKey := fmt.Sprintf("%s:%s", authBanKey, id)

	ttl, err := s.rdb.TTL(ctx, banKey).Result()
	if err != nil && err != redis.Nil {
		s.log.Error(err, "failed to get auth ban ttl in redis")
		return 0, err
	}
	if ttl > 0 {
		return ttl, nil
	}
	return 0, nil
}

// IncrementAttempts: dipanggil kalau login/refresh gagal
func (s *Security) IncrementAttempts(ctx context.Context, id string) (time.Duration, error) {
	attemptsKey := fmt.Sprintf("%s:%s", authAttemptsKey, id)
	banKey := fmt.Sprintf("%s:%s", authBanKey, id)

	// increment attempts
	attempts, err := s.rdb.Incr(ctx, attemptsKey).Result()
	if err != nil {
		s.log.Error(err, "failed to increment auth attempts in redis")
		return 0, err
	}

	// TTL panjang buat counter (misalnya 1 jam sejak first fail)
	if attempts == 1 {
		_ = s.rdb.Expire(ctx, attemptsKey, 1*time.Hour).Err()
	}

	// kalau attempts > 3 → hitung delay ban
	if attempts > 3 {
		delay := time.Duration(5*(1<<(attempts-3))) * time.Second
		if delay > 1*time.Hour {
			delay = 1 * time.Hour
		}

		// set ban key dengan TTL delay
		if err := s.rdb.SetEx(ctx, banKey, "1", delay).Err(); err != nil {
			s.log.Error(err, "failed to set auth ban key in redis")
			return 0, err
		}
		return delay, nil
	}

	return 0, nil
}

// ResetAuthAttempts: dipanggil kalau login/refresh sukses
func (s *Security) ResetAuthAttempts(ctx context.Context, id string) error {
	attemptsKey := fmt.Sprintf("%s:%s", authAttemptsKey, id)
	banKey := fmt.Sprintf("%s:%s", authBanKey, id)

	_, err := s.rdb.Del(ctx, attemptsKey, banKey).Result()
	if err != nil {
		s.log.Error(err, "failed to reset auth attempts in redis")
		return err
	}
	return nil
}
