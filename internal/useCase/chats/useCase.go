package chat

import (
	"errors"
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/repository/storage/postgres"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"

	"github.com/google/uuid"
)

type chatsUC struct {
	cfg                  *config.Config
	log                  *logger.ZapLogger
	chatRepo             postgres.ChatRepository
	chatParticipantsRepo postgres.ChatParticipantRepository
}

func New(cfg *config.Config, chatRepo postgres.ChatRepository, chatParticipantsRepo postgres.ChatParticipantRepository, log *logger.ZapLogger) *chatsUC {
	return &chatsUC{cfg, log, chatRepo, chatParticipantsRepo}
}

func (uc *chatsUC) CreateChatWithParticipants(chatWithParticipants *models.ChatWithParticipants) (*models.ChatWithParticipants, error) {
	if chatWithParticipants.Chat.ChatType == models.PrivateChat && len(chatWithParticipants.Participants) != 2 {
		return nil, errors.New(e.ErrAllowedParticipants)
	}

	existingChat, err := uc.chatParticipantsRepo.GetChatByParticipants(chatWithParticipants.Participants)
	if err != nil {
		uc.log.Error("Error checking for existing chat:", err)
		return nil, err
	}

	if existingChat != nil && existingChat.ChatType == models.PrivateChat {
		return nil, errors.New(e.ErrChatExists)
	}

	createdChat, err := uc.chatRepo.CreateChat(&chatWithParticipants.Chat)
	if err != nil {
		uc.log.Error("Error in CreateChat:", err)
		return nil, err
	}

	createdParticipants := make([]models.ChatParticipant, 0)

	for _, participant := range chatWithParticipants.Participants {
		participant.ChatID = createdChat.ChatID
		err := uc.chatParticipantsRepo.CreateChatParticipant(&participant)
		if err != nil {
			return nil, err
		}
		createdParticipants = append(createdParticipants, participant)
	}

	chatWithParticipants.Participants = createdParticipants

	return chatWithParticipants, nil
}

func (uc *chatsUC) ListChatsWithParticipants(userID uuid.UUID) ([]*models.ChatWithParticipants, error) {
	// Retrieve the list of chat IDs where the current user is a participant
	chats, err := uc.chatParticipantsRepo.GetChatsByUserID(userID)
	if err != nil {
		return nil, err
	}
	chatList := make([]*models.ChatWithParticipants, 0)

	for _, chat := range chats {
		chat, err := uc.chatRepo.GetChatByID(chat.ChatID)
		if err != nil {
			return nil, err
		}

		participants, err := uc.chatParticipantsRepo.GetParticipantsByChatID(chat.ChatID)
		if err != nil {
			return nil, err
		}

		chatWithParticipants := &models.ChatWithParticipants{
			Chat:         *chat,
			Participants: participants,
		}

		chatList = append(chatList, chatWithParticipants)
	}

	return chatList, nil
}

func (uc *chatsUC) GetChatByID(chatID uuid.UUID) (*models.ChatWithParticipants, error) {
	chat, err := uc.chatRepo.GetChatByID(chatID)
	if err != nil {
		return nil, err
	}

	participants, err := uc.chatParticipantsRepo.GetParticipantsByChatID(chat.ChatID)
	if err != nil {
		return nil, err
	}

	chatWithParticipants := &models.ChatWithParticipants{
		Chat:         *chat,
		Participants: participants,
	}

	return chatWithParticipants, nil
}

func (uc *chatsUC) DeleteChat(chatID uuid.UUID, requestingUserID uuid.UUID) error {
	// Get the participants of the chat
	participants, err := uc.chatParticipantsRepo.GetParticipantsByChatID(chatID)
	if err != nil {
		return err
	}

	// Check if the requesting user is a participant in the chat
	isParticipant := false
	for _, participant := range participants {
		if participant.UserID == requestingUserID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return errors.New(e.ErrForbidden)
	}

	// Delete the participants from the chat
	err = uc.chatParticipantsRepo.DeleteParticipantsByChatID(chatID)
	if err != nil {
		return err
	}

	// Delete the chat
	err = uc.chatRepo.DeleteChat(chatID)
	if err != nil {
		return err
	}

	return nil
}
