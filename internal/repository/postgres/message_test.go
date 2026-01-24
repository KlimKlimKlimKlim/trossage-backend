package postgres

import (
	"context"
	"strconv"

	"github.com/jackc/pgerrcode"
)

//nolint:gosmopolitan // Testing Unicode support in different scripts
func (s *RepositoryTestSuite) TestCreateMessage() {
	testCases := []struct {
		name         string
		loginSuffix  string
		messageText  string
		expectedErr  string
		chatExists   bool
		senderExists bool
		addToChat    bool
	}{
		{
			name:         "success - creates message successfully",
			loginSuffix:  "msg1",
			messageText:  "Hello, world!",
			chatExists:   true,
			senderExists: true,
			addToChat:    true,
			expectedErr:  "",
		},
		{
			name:         "success - creates message with unicode text",
			loginSuffix:  "msg2",
			messageText:  "Привет, мир! 你好世界",
			chatExists:   true,
			senderExists: true,
			addToChat:    true,
			expectedErr:  "",
		},
		{
			name:         "success - creates empty message",
			loginSuffix:  "msg3",
			messageText:  "",
			chatExists:   true,
			senderExists: true,
			addToChat:    true,
			expectedErr:  "",
		},
		{
			name:         "error - chat does not exist",
			loginSuffix:  "msg4",
			messageText:  "Message",
			chatExists:   false,
			senderExists: true,
			addToChat:    false,
			expectedErr:  pgerrcode.ForeignKeyViolation,
		},
		{
			name:         "error - sender does not exist",
			loginSuffix:  "msg5",
			messageText:  "Message",
			chatExists:   true,
			senderExists: false,
			addToChat:    false,
			expectedErr:  pgerrcode.ForeignKeyViolation,
		},
		{
			name:         "success - sender not in chat members",
			loginSuffix:  "msg6",
			messageText:  "Message from non-member",
			chatExists:   true,
			senderExists: true,
			addToChat:    false,
			expectedErr:  "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var chatID, senderID int64

			if tc.chatExists {
				chat := s.createTestChat(ctx)
				chatID = chat.ID
			} else {
				chatID = 999999
			}

			if tc.senderExists {
				sender := s.createTestUser(ctx, "sender_"+tc.loginSuffix)
				senderID = sender.ID

				if tc.addToChat && tc.chatExists {
					err := s.repo.InsertChatParticipants(ctx, chatID, senderID)
					s.Require().NoError(err)
				}
			} else {
				senderID = 999999
			}

			// Test.
			result, err := s.repo.CreateMessage(ctx, chatID, senderID, tc.messageText)

			// Check.
			if tc.expectedErr != "" {
				s.assertPgError(err, tc.expectedErr)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.NotZero(result.ID)
			s.Equal(chatID, result.ChatID)
			s.Equal(senderID, result.SenderID)
			s.Equal(tc.messageText, result.Text)
			s.NotZero(result.CreatedAt)
		})
	}
}

func (s *RepositoryTestSuite) TestSelectMessages() {
	testCases := []struct {
		name          string
		loginSuffix   string
		setupMessages int
		limit         int
		offset        int
		expectedCount int
		chatExists    bool
	}{
		{
			name:          "success - returns messages ordered by created_at desc",
			loginSuffix:   "selmsg1",
			setupMessages: 5,
			limit:         10,
			offset:        0,
			expectedCount: 5,
			chatExists:    true,
		},
		{
			name:          "success - pagination with limit",
			loginSuffix:   "selmsg2",
			setupMessages: 10,
			limit:         3,
			offset:        0,
			expectedCount: 3,
			chatExists:    true,
		},
		{
			name:          "success - pagination with offset",
			loginSuffix:   "selmsg3",
			setupMessages: 10,
			limit:         3,
			offset:        5,
			expectedCount: 3,
			chatExists:    true,
		},
		{
			name:          "success - empty result for chat without messages",
			loginSuffix:   "selmsg4",
			setupMessages: 0,
			limit:         10,
			offset:        0,
			expectedCount: 0,
			chatExists:    true,
		},
		{
			name:          "success - offset beyond results",
			loginSuffix:   "selmsg5",
			setupMessages: 3,
			limit:         10,
			offset:        10,
			expectedCount: 0,
			chatExists:    true,
		},
		{
			name:          "success - non-existent chat returns empty",
			loginSuffix:   "selmsg6",
			setupMessages: 0,
			limit:         10,
			offset:        0,
			expectedCount: 0,
			chatExists:    false,
		},
		{
			name:          "success - limit zero returns empty",
			loginSuffix:   "selmsg7",
			setupMessages: 5,
			limit:         0,
			offset:        0,
			expectedCount: 0,
			chatExists:    true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var chatID int64

			if tc.chatExists {
				chat := s.createTestChat(ctx)
				chatID = chat.ID
			} else {
				chatID = 999999
			}

			var senderID int64

			if tc.setupMessages > 0 {
				sender := s.createTestUser(ctx, "sender_"+tc.loginSuffix)
				senderID = sender.ID

				err := s.repo.InsertChatParticipants(ctx, chatID, senderID)
				s.Require().NoError(err)

				for i := range tc.setupMessages {
					_, err = s.repo.CreateMessage(ctx, chatID, senderID, "Message "+strconv.Itoa(i))
					s.Require().NoError(err)
				}
			}

			// Test.
			result, err := s.repo.SelectMessages(ctx, chatID, tc.limit, tc.offset)

			// Check.
			s.Require().NoError(err)
			s.Len(result, tc.expectedCount)

			if tc.expectedCount > 0 {
				for i, msg := range result {
					s.NotZero(msg.ID)
					s.Equal(chatID, msg.ChatID)
					s.Equal(senderID, msg.SenderID)
					s.NotEmpty(msg.Text)
					s.NotZero(msg.CreatedAt)

					if i > 0 {
						s.True(result[i-1].CreatedAt.After(msg.CreatedAt) || result[i-1].CreatedAt.Equal(msg.CreatedAt),
							"messages should be ordered by created_at DESC")
					}
				}

				if tc.offset == 0 && tc.setupMessages > 0 {
					expectedLastText := "Message " + strconv.Itoa(tc.setupMessages-1)
					s.Equal(expectedLastText, result[0].Text, "first message should be the last created")
				}
			}
		})
	}
}

func (s *RepositoryTestSuite) TestCountChatMessages() {
	testCases := []struct {
		name          string
		loginSuffix   string
		setupMessages int
		expectedCount int
		chatExists    bool
	}{
		{
			name:          "success - counts chat messages",
			loginSuffix:   "countmsg1",
			setupMessages: 7,
			expectedCount: 7,
			chatExists:    true,
		},
		{
			name:          "success - returns zero for chat without messages",
			loginSuffix:   "countmsg2",
			setupMessages: 0,
			expectedCount: 0,
			chatExists:    true,
		},
		{
			name:          "success - non-existent chat returns zero",
			loginSuffix:   "countmsg3",
			setupMessages: 0,
			expectedCount: 0,
			chatExists:    false,
		},
		{
			name:          "success - counts single message",
			loginSuffix:   "countmsg4",
			setupMessages: 1,
			expectedCount: 1,
			chatExists:    true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var chatID int64

			if tc.chatExists {
				chat := s.createTestChat(ctx)
				chatID = chat.ID
			} else {
				chatID = 999999
			}

			if tc.setupMessages > 0 {
				sender := s.createTestUser(ctx, "sender_"+tc.loginSuffix)

				err := s.repo.InsertChatParticipants(ctx, chatID, sender.ID)
				s.Require().NoError(err)

				for i := range tc.setupMessages {
					_, err = s.repo.CreateMessage(ctx, chatID, sender.ID, "Message "+strconv.Itoa(i))
					s.Require().NoError(err)
				}
			}

			// Test.
			result, err := s.repo.CountChatMessages(ctx, chatID)

			// Check.
			s.Require().NoError(err)
			s.Equal(tc.expectedCount, result)
		})
	}
}
