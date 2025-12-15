package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"video_conference/internal/domain"
	"video_conference/internal/repository"
	"video_conference/pkg/logger"
)

type ChatService interface {
	SendMessage(ctx context.Context, roomID uuid.UUID, participantID uuid.UUID, content string) (*domain.ChatMessage, error)
	GetMessages(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*domain.ChatMessage, error)
	EditMessage(ctx context.Context, messageID int64, participantID uuid.UUID, content string) (*domain.ChatMessage, error)
	DeleteMessage(ctx context.Context, messageID int64, participantID uuid.UUID) error
}

type chatService struct {
	chatRepo repository.ChatRepository
	roomRepo repository.RoomRepository
	auditRepo repository.AuditRepository
	log      logger.Logger
}

func NewChatService(chatRepo repository.ChatRepository, roomRepo repository.RoomRepository, auditRepo repository.AuditRepository, log logger.Logger) ChatService {
	return &chatService{
		chatRepo:  chatRepo,
		roomRepo:  roomRepo,
		auditRepo: auditRepo,
		log:       log,
	}
}

func (s *chatService) SendMessage(ctx context.Context, roomID uuid.UUID, participantID uuid.UUID, content string) (*domain.ChatMessage, error) {
	// Проверка существования комнаты
	_, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	message := &domain.ChatMessage{
		RoomID:              roomID,
		SenderParticipantID: &participantID,
		MessageType:         domain.MessageTypeUser,
		Content:             content,
		CreatedAt:           time.Now(),
	}

	if err := s.chatRepo.CreateMessage(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *chatService) GetMessages(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*domain.ChatMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.chatRepo.GetMessages(ctx, roomID, limit, offset)
}

func (s *chatService) EditMessage(ctx context.Context, messageID int64, participantID uuid.UUID, content string) (*domain.ChatMessage, error) {
	message, err := s.chatRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if message.SenderParticipantID == nil || *message.SenderParticipantID != participantID {
		return nil, errors.New("only sender can edit message")
	}

	message.Content = content
	if err := s.chatRepo.UpdateMessage(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *chatService) DeleteMessage(ctx context.Context, messageID int64, participantID uuid.UUID) error {
	message, err := s.chatRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}

	if message.SenderParticipantID == nil || *message.SenderParticipantID != participantID {
		return errors.New("only sender can delete message")
	}

	return s.chatRepo.DeleteMessage(ctx, messageID, participantID)
}

