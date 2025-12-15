package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/livekit/protocol/auth"
	"video_conference/internal/config"
	"video_conference/internal/repository"
	"video_conference/pkg/logger"
)

type MediaService interface {
	GetToken(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, displayName string) (string, error)
}

type mediaService struct {
	roomRepo repository.RoomRepository
	cfg      config.LiveKitConfig
	log      logger.Logger
}

func NewMediaService(roomRepo repository.RoomRepository, cfg config.LiveKitConfig, log logger.Logger) MediaService {
	return &mediaService{
		roomRepo: roomRepo,
		cfg:      cfg,
		log:      log,
	}
}

func (s *mediaService) GetToken(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, displayName string) (string, error) {
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return "", errors.New("room not found")
	}

	at := auth.NewAccessToken(s.cfg.APIKey, s.cfg.APISecret)
	canPublish := true
	canSubscribe := true
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         room.LiveKitRoomName,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	}

	at.AddGrant(grant).
		SetIdentity(userID.String()).
		SetName(displayName).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		s.log.Error("Failed to generate LiveKit token", "error", err)
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

