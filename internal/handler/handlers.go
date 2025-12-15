package handler

import (
	"video_conference/internal/service"
	"video_conference/pkg/logger"
)

type Handlers struct {
	Health     *HealthHandler
	Auth       *AuthHandler
	User       *UserHandler
	Room       *RoomHandler
	WaitingRoom *WaitingRoomHandler
	Chat       *ChatHandler
	Media      *MediaHandler
	Stats      *StatsHandler
	WebSocket  *WebSocketHandler
}

func NewHandlers(services *service.Services, log logger.Logger) *Handlers {
	return &Handlers{
		Health:      NewHealthHandler(),
		Auth:       NewAuthHandler(services.Auth, log),
		User:       NewUserHandler(services.User, log),
		Room:       NewRoomHandler(services.Room, log),
		WaitingRoom: NewWaitingRoomHandler(services.Room, log),
		Chat:       NewChatHandler(services.Chat, log),
		Media:      NewMediaHandler(services.Media, log),
		Stats:      NewStatsHandler(services.Stats, log),
		WebSocket:  NewWebSocketHandler(services.Chat, log),
	}
}

