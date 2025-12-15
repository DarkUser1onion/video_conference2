package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"video_conference/pkg/logger"
)

type Repositories struct {
	User      UserRepository
	Room      RoomRepository
	Chat      ChatRepository
	Stats     StatsRepository
	Audit     AuditRepository
	RateLimit RateLimitRepository
}

func NewRepositories(db *pgxpool.Pool, redis *redis.Client, log logger.Logger) *Repositories {
	return &Repositories{
		User:      NewUserRepository(db, log),
		Room:      NewRoomRepository(db, log),
		Chat:      NewChatRepository(db, log),
		Stats:     NewStatsRepository(db, log),
		Audit:     NewAuditRepository(db, log),
		RateLimit: NewRateLimitRepository(redis, log),
	}
}

