package service

import (
	"video_conference/internal/config"
	"video_conference/internal/repository"
	"video_conference/pkg/logger"
)

type Services struct {
	Auth          AuthService
	User          UserService
	Room          RoomService
	Chat          ChatService
	Media         MediaService
	Stats         StatsService
	RateLimit     RateLimitService
	Audit         AuditService
	ScreenCapture ScreenCaptureService
	AudioCapture  AudioCaptureService
	WebRTC        WebRTCService
}

func NewServices(repos *repository.Repositories, cfg *config.Config, log logger.Logger) *Services {
	return &Services{
		Auth:          NewAuthService(repos.User, cfg.JWT, log),
		User:          NewUserService(repos.User, repos.Audit, log),
		Room:          NewRoomService(repos.Room, repos.Audit, cfg, log),
		Chat:          NewChatService(repos.Chat, repos.Room, repos.Audit, log),
		Media:         NewMediaService(repos.Room, cfg.LiveKit, log),
		Stats:         NewStatsService(repos.Stats, log),
		RateLimit:     NewRateLimitService(repos.RateLimit, log),
		Audit:         NewAuditService(repos.Audit, log),
		ScreenCapture: NewScreenCaptureService(log),
		AudioCapture:  NewAudioCaptureService(log),
		WebRTC:        NewWebRTCService(log),
	}
}

// NewServicesForScreenShare создает только сервисы для screen share без зависимостей от БД
func NewServicesForScreenShare(cfg *config.Config, log logger.Logger) *Services {
	return &Services{
		ScreenCapture: NewScreenCaptureService(log),
		AudioCapture:  NewAudioCaptureService(log),
		WebRTC:        NewWebRTCService(log),
		// Остальные сервисы будут nil, но это не проблема для screen share
	}
}
