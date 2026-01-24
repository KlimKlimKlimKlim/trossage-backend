package postgres

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/jackc/pgerrcode"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

func (s *RepositoryTestSuite) TestSelectChatBetweenUsers() {
	testCases := []struct {
		expectedErr error
		name        string
		loginSuffix string
		setupUsers  int
		userID1Idx  int
		userID2Idx  int
		createChat  bool
		wrongUserID bool
	}{
		{
			name:        "success - finds chat between two users",
			loginSuffix: "chat1",
			setupUsers:  2,
			createChat:  true,
			userID1Idx:  0,
			userID2Idx:  1,
			wrongUserID: false,
			expectedErr: nil,
		},
		{
			name:        "error - chat not found between users",
			loginSuffix: "chat2",
			setupUsers:  2,
			createChat:  false,
			userID1Idx:  0,
			userID2Idx:  1,
			wrongUserID: false,
			expectedErr: derrors.ErrChatNotFound,
		},
		{
			name:        "error - chat not found with non-existent user",
			loginSuffix: "chat3",
			setupUsers:  2,
			createChat:  true,
			userID1Idx:  0,
			userID2Idx:  1,
			wrongUserID: true,
			expectedErr: derrors.ErrChatNotFound,
		},
		{
			name:        "success - finds chat with reversed user order",
			loginSuffix: "chat4",
			setupUsers:  2,
			createChat:  true,
			userID1Idx:  1,
			userID2Idx:  0,
			wrongUserID: false,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			users := make([]model.AuthUser, 0, tc.setupUsers)
			for i := range tc.setupUsers {
				user := s.createTestUser(ctx, fmt.Sprintf("%s_%d", tc.loginSuffix, i))
				users = append(users, user)
			}

			var expectedChatID int64

			if tc.createChat {
				chat := s.createTestChat(ctx)
				expectedChatID = chat.ID

				err := s.repo.InsertChatParticipants(ctx, chat.ID, users[0].ID, users[1].ID)
				s.Require().NoError(err)
			}

			queryUserID1 := users[tc.userID1Idx].ID

			queryUserID2 := users[tc.userID2Idx].ID
			if tc.wrongUserID {
				queryUserID2 = 999999
			}

			// Test.
			chatID, err := s.repo.SelectChatBetweenUsers(ctx, queryUserID1, queryUserID2)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(chatID)

				return
			}

			s.Require().NoError(err)
			s.Equal(expectedChatID, chatID)
		})
	}
}

func (s *RepositoryTestSuite) TestInsertChat() {
	testCases := []struct {
		name        string
		chatType    model.ChatType
		expectedErr string
	}{
		{
			name:        "success - creates private chat",
			chatType:    model.ChatTypePrivate,
			expectedErr: "",
		},
		{
			name:        "error - invalid chat type",
			chatType:    model.ChatType("invalid"),
			expectedErr: pgerrcode.InvalidTextRepresentation,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Test.
			result, err := s.repo.InsertChat(ctx, tc.chatType)

			// Check.
			if tc.expectedErr != "" {
				s.assertPgError(err, tc.expectedErr)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.NotZero(result.ID)
			s.Equal(tc.chatType, result.Type)
			s.NotZero(result.CreatedAt)
			s.NotZero(result.UpdatedAt)
			s.Empty(result.Title)
			s.True(result.DeletedAt.IsZero())
		})
	}
}

func (s *RepositoryTestSuite) TestInsertChatParticipants() {
	testCases := []struct {
		expectedErr error
		name        string
		loginSuffix string
		pgCode      string
		userCount   int
		chatExists  bool
		usersExist  bool
		duplicate   bool
	}{
		{
			name:        "success - inserts single participant",
			loginSuffix: "part1",
			userCount:   1,
			chatExists:  true,
			usersExist:  true,
			expectedErr: nil,
		},
		{
			name:        "success - inserts multiple participants",
			loginSuffix: "part2",
			userCount:   3,
			chatExists:  true,
			usersExist:  true,
			expectedErr: nil,
		},
		{
			name:        "error - empty user IDs",
			loginSuffix: "part3",
			userCount:   0,
			chatExists:  true,
			usersExist:  false,
			expectedErr: derrors.ErrEmptyInput,
		},
		{
			name:        "error - chat does not exist",
			loginSuffix: "part4",
			userCount:   1,
			chatExists:  false,
			usersExist:  true,
			pgCode:      pgerrcode.ForeignKeyViolation,
		},
		{
			name:        "error - user does not exist",
			loginSuffix: "part5",
			userCount:   1,
			chatExists:  true,
			usersExist:  false,
			pgCode:      pgerrcode.ForeignKeyViolation,
		},
		{
			name:        "error - duplicate participant",
			loginSuffix: "part6",
			userCount:   1,
			chatExists:  true,
			usersExist:  true,
			duplicate:   true,
			pgCode:      pgerrcode.UniqueViolation,
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

			userIDs := make([]int64, 0, tc.userCount)
			if tc.usersExist {
				for i := range tc.userCount {
					user := s.createTestUser(ctx, tc.loginSuffix+"_"+strconv.Itoa(i))
					userIDs = append(userIDs, user.ID)
				}
			} else if tc.userCount > 0 {
				userIDs = []int64{999999}
			}

			if tc.duplicate && len(userIDs) > 0 {
				err := s.repo.InsertChatParticipants(ctx, chatID, userIDs[0])
				s.Require().NoError(err)
			}

			// Test.
			err := s.repo.InsertChatParticipants(ctx, chatID, userIDs...)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				return
			}

			if tc.pgCode != "" {
				s.assertPgError(err, tc.pgCode)

				return
			}

			s.Require().NoError(err)
		})
	}
}

func (s *RepositoryTestSuite) TestSelectUserChats() {
	testCases := []struct {
		name                   string
		loginSuffix            string
		withMessages           []int
		setupChats             int
		limit                  int
		offset                 int
		expectedCount          int
		userExists             bool
		expectChatsWithoutMsgs bool
	}{
		{
			name:                   "success - returns chats ordered by last message",
			loginSuffix:            "chats1",
			setupChats:             3,
			withMessages:           []int{0, 1, 2},
			limit:                  10,
			offset:                 0,
			expectedCount:          3,
			userExists:             true,
			expectChatsWithoutMsgs: false,
		},
		{
			name:                   "success - returns chats without messages last",
			loginSuffix:            "chats2",
			setupChats:             3,
			withMessages:           []int{0},
			limit:                  10,
			offset:                 0,
			expectedCount:          3,
			userExists:             true,
			expectChatsWithoutMsgs: true,
		},
		{
			name:                   "success - pagination with limit",
			loginSuffix:            "chats3",
			setupChats:             5,
			withMessages:           []int{0, 1, 2, 3, 4},
			limit:                  2,
			offset:                 0,
			expectedCount:          2,
			userExists:             true,
			expectChatsWithoutMsgs: false,
		},
		{
			name:                   "success - pagination with offset",
			loginSuffix:            "chats4",
			setupChats:             5,
			withMessages:           []int{0, 1, 2, 3, 4},
			limit:                  2,
			offset:                 2,
			expectedCount:          2,
			userExists:             true,
			expectChatsWithoutMsgs: false,
		},
		{
			name:                   "success - empty result for user with no chats",
			loginSuffix:            "chats5",
			setupChats:             0,
			withMessages:           []int{},
			limit:                  10,
			offset:                 0,
			expectedCount:          0,
			userExists:             true,
			expectChatsWithoutMsgs: false,
		},
		{
			name:                   "success - offset beyond results",
			loginSuffix:            "chats6",
			setupChats:             2,
			withMessages:           []int{0, 1},
			limit:                  10,
			offset:                 5,
			expectedCount:          0,
			userExists:             true,
			expectChatsWithoutMsgs: false,
		},
		{
			name:                   "success - user does not exist",
			loginSuffix:            "chats7",
			setupChats:             0,
			withMessages:           []int{},
			limit:                  10,
			offset:                 0,
			expectedCount:          0,
			userExists:             false,
			expectChatsWithoutMsgs: false,
		},
		{
			name:                   "success - limit zero returns empty",
			loginSuffix:            "chats8",
			setupChats:             3,
			withMessages:           []int{0, 1, 2},
			limit:                  0,
			offset:                 0,
			expectedCount:          0,
			userExists:             true,
			expectChatsWithoutMsgs: false,
		},
		{
			name:                   "success - very large limit returns all",
			loginSuffix:            "chats9",
			setupChats:             3,
			withMessages:           []int{0, 1, 2},
			limit:                  1000,
			offset:                 0,
			expectedCount:          3,
			userExists:             true,
			expectChatsWithoutMsgs: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var userID int64

			if !tc.userExists {
				userID = 999999
			} else {
				mainUser := s.createTestUser(ctx, "chatuser_"+tc.loginSuffix)
				userID = mainUser.ID

				for i := range tc.setupChats {
					chat := s.createTestChat(ctx)
					otherUser := s.createTestUser(ctx, "other_"+tc.loginSuffix+"_"+strconv.Itoa(i))

					err := s.repo.InsertChatParticipants(ctx, chat.ID, userID, otherUser.ID)
					s.Require().NoError(err)

					hasMessage := slices.Contains(tc.withMessages, i)

					if hasMessage {
						_, err = s.repo.CreateMessage(ctx, chat.ID, otherUser.ID, "Message "+strconv.Itoa(i))
						s.Require().NoError(err)
					}
				}
			}

			// Test.
			result, err := s.repo.SelectUserChats(ctx, userID, tc.limit, tc.offset)

			// Check.
			s.Require().NoError(err)
			s.Len(result, tc.expectedCount)

			if tc.expectedCount > 0 {
				for _, chat := range result {
					s.NotZero(chat.ID)
					s.Equal(model.ChatTypePrivate, chat.Type)
					s.NotZero(chat.OtherUser.ID)
					s.NotEmpty(chat.OtherUser.Login)
				}

				if len(tc.withMessages) > 0 && tc.offset == 0 {
					s.False(result[0].LastMessage.IsEmpty(), "first chat should have a message")
				}

				if tc.expectChatsWithoutMsgs {
					hasMessageFirst := !result[0].LastMessage.IsEmpty()
					s.True(hasMessageFirst, "chat with message should be first")

					for i := 1; i < len(result); i++ {
						s.True(result[i].LastMessage.IsEmpty(), "chats without messages should be last")
					}
				}
			}
		})
	}
}

func (s *RepositoryTestSuite) TestCountUserChats() {
	testCases := []struct {
		name          string
		loginSuffix   string
		setupChats    int
		expectedCount int
		userExists    bool
	}{
		{
			name:          "success - counts user chats",
			loginSuffix:   "count1",
			setupChats:    5,
			expectedCount: 5,
			userExists:    true,
		},
		{
			name:          "success - returns zero for user with no chats",
			loginSuffix:   "count2",
			setupChats:    0,
			expectedCount: 0,
			userExists:    true,
		},
		{
			name:          "success - user does not exist",
			loginSuffix:   "count3",
			setupChats:    0,
			expectedCount: 0,
			userExists:    false,
		},
		{
			name:          "success - counts single chat",
			loginSuffix:   "count4",
			setupChats:    1,
			expectedCount: 1,
			userExists:    true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var userID int64

			if tc.userExists {
				mainUser := s.createTestUser(ctx, "countuser_"+tc.loginSuffix)
				userID = mainUser.ID

				for i := range tc.setupChats {
					chat := s.createTestChat(ctx)
					otherUser := s.createTestUser(ctx, "other_"+tc.loginSuffix+"_"+strconv.Itoa(i))

					err := s.repo.InsertChatParticipants(ctx, chat.ID, userID, otherUser.ID)
					s.Require().NoError(err)
				}
			} else {
				userID = 999999
			}

			// Test.
			result, err := s.repo.CountUserChats(ctx, userID)

			// Check.
			s.Require().NoError(err)
			s.Equal(tc.expectedCount, result)
		})
	}
}

func (s *RepositoryTestSuite) TestIsUserMember() {
	testCases := []struct {
		name           string
		loginSuffix    string
		expectedResult bool
		chatExists     bool
		userExists     bool
		addParticipant bool
	}{
		{
			name:           "success - user is member of chat",
			loginSuffix:    "member1",
			chatExists:     true,
			userExists:     true,
			addParticipant: true,
			expectedResult: true,
		},
		{
			name:           "success - user is not member of chat",
			loginSuffix:    "member2",
			chatExists:     true,
			userExists:     true,
			addParticipant: false,
			expectedResult: false,
		},
		{
			name:           "success - chat does not exist",
			loginSuffix:    "member3",
			chatExists:     false,
			userExists:     true,
			addParticipant: false,
			expectedResult: false,
		},
		{
			name:           "success - user does not exist",
			loginSuffix:    "member4",
			chatExists:     true,
			userExists:     false,
			addParticipant: false,
			expectedResult: false,
		},
		{
			name:           "success - both do not exist",
			loginSuffix:    "member5",
			chatExists:     false,
			userExists:     false,
			addParticipant: false,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var chatID, userID int64

			if tc.chatExists {
				chat := s.createTestChat(ctx)
				chatID = chat.ID
			} else {
				chatID = 999999
			}

			if tc.userExists {
				user := s.createTestUser(ctx, "memberuser_"+tc.loginSuffix)
				userID = user.ID
			} else {
				userID = 999999
			}

			if tc.addParticipant && tc.chatExists && tc.userExists {
				err := s.repo.InsertChatParticipants(ctx, chatID, userID)
				s.Require().NoError(err)
			}

			// Test.
			result, err := s.repo.IsUserMember(ctx, chatID, userID)

			// Check.
			s.Require().NoError(err)
			s.Equal(tc.expectedResult, result)
		})
	}
}

func (s *RepositoryTestSuite) TestSelectChatByID() {
	testCases := []struct {
		expectedErr error
		name        string
		loginSuffix string
		chatExists  bool
		deleteChat  bool
	}{
		{
			name:        "success - returns existing chat",
			loginSuffix: "selectchat1",
			chatExists:  true,
			deleteChat:  false,
			expectedErr: nil,
		},
		{
			name:        "error - chat not found",
			loginSuffix: "selectchat2",
			chatExists:  false,
			deleteChat:  false,
			expectedErr: derrors.ErrChatNotFound,
		},
		{
			name:        "error - chat is deleted",
			loginSuffix: "selectchat3",
			chatExists:  true,
			deleteChat:  true,
			expectedErr: derrors.ErrChatNotFound,
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

				if tc.deleteChat {
					_, err := s.pool.Exec(ctx, "UPDATE chats SET deleted_at = NOW() WHERE id = $1", chatID)
					s.Require().NoError(err)
				}
			} else {
				chatID = 999999
			}

			// Test.
			result, err := s.repo.SelectChatByID(ctx, chatID)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.Equal(chatID, result.ID)
			s.Equal(model.ChatTypePrivate, result.Type)
			s.NotZero(result.CreatedAt)
			s.NotZero(result.UpdatedAt)
			s.True(result.DeletedAt.IsZero())
		})
	}
}

func (s *RepositoryTestSuite) TestSelectChatMembers() {
	testCases := []struct {
		name          string
		loginSuffix   string
		setupMembers  int
		expectedCount int
		chatExists    bool
	}{
		{
			name:          "success - returns all chat members",
			loginSuffix:   "members1",
			setupMembers:  3,
			expectedCount: 3,
			chatExists:    true,
		},
		{
			name:          "success - returns single member",
			loginSuffix:   "members2",
			setupMembers:  1,
			expectedCount: 1,
			chatExists:    true,
		},
		{
			name:          "success - returns empty for chat without members",
			loginSuffix:   "members3",
			setupMembers:  0,
			expectedCount: 0,
			chatExists:    true,
		},
		{
			name:          "success - returns empty for non-existent chat",
			loginSuffix:   "members4",
			setupMembers:  0,
			expectedCount: 0,
			chatExists:    false,
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

			expectedUserIDs := make([]int64, 0, tc.setupMembers)

			for i := range tc.setupMembers {
				user := s.createTestUser(ctx, "chatmember_"+tc.loginSuffix+"_"+strconv.Itoa(i))
				expectedUserIDs = append(expectedUserIDs, user.ID)
			}

			if tc.setupMembers > 0 && tc.chatExists {
				err := s.repo.InsertChatParticipants(ctx, chatID, expectedUserIDs...)
				s.Require().NoError(err)
			}

			// Test.
			result, err := s.repo.SelectChatMembers(ctx, chatID)

			// Check.
			s.Require().NoError(err)
			s.Len(result, tc.expectedCount)

			if tc.expectedCount > 0 {
				s.ElementsMatch(expectedUserIDs, result)
			}
		})
	}
}
