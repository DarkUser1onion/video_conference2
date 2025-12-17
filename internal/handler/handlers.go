package handler

import (
	"video_conference/internal/service"
	"video_conference/pkg/logger"
)

type Handlers struct {
	Health      *HealthHandler
	Auth        *AuthHandler
	User        *UserHandler
	Room        *RoomHandler
	WaitingRoom *WaitingRoomHandler
	Chat        *ChatHandler
	Media       *MediaHandler
	Stats       *StatsHandler
	WebSocket   *WebSocketHandler
	ScreenShare *ScreenShareHandler
}

func NewHandlers(services *service.Services, log logger.Logger) *Handlers {
	return &Handlers{
		Health:      NewHealthHandler(),
		Auth:        NewAuthHandler(services.Auth, log),
		User:        NewUserHandler(services.User, log),
		Room:        NewRoomHandler(services.Room, log),
		WaitingRoom: NewWaitingRoomHandler(services.Room, log),
		Chat:        NewChatHandler(services.Chat, log),
		Media:       NewMediaHandler(services.Media, log),
		Stats:       NewStatsHandler(services.Stats, log),
		WebSocket:   NewWebSocketHandler(services.Chat, log),
		ScreenShare: NewScreenShareHandler(services.ScreenCapture, services.AudioCapture, services.WebRTC, log),
	}
}

// NewHandlersForScreenShare создает только handlers для screen share
func NewHandlersForScreenShare(services *service.Services, log logger.Logger) *Handlers {
	return &Handlers{
		Health:      NewHealthHandler(),
		ScreenShare: NewScreenShareHandler(services.ScreenCapture, services.AudioCapture, services.WebRTC, log),
		// Остальные handlers будут nil, но это не проблема для screen share endpoints
	}
}
