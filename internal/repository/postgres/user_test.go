package postgres

import (
	"context"
	"errors"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

//nolint:gosmopolitan // Testing Unicode support in different scripts
func (s *RepositoryTestSuite) TestInsertUser() {
	testCases := []struct {
		expectedErr error
		name        string
		user        model.AuthUser
	}{
		{
			name: "success - valid user with all fields",
			user: model.AuthUser{
				User: model.User{
					Login:       "validuser",
					DisplayName: "Valid User",
				},
				PasswordHash: "hash123",
			},
			expectedErr: nil,
		},
		{
			name: "success - display name with cyrillic",
			user: model.AuthUser{
				User: model.User{
					Login:       "cyrillicuser",
					DisplayName: "Пользователь Иванов",
				},
				PasswordHash: "hash_cyrillic",
			},
			expectedErr: nil,
		},
		{
			name: "success - display name with chinese",
			user: model.AuthUser{
				User: model.User{
					Login:       "chineseuser",
					DisplayName: "用户王小明",
				},
				PasswordHash: "hash_chinese",
			},
			expectedErr: nil,
		},
		{
			name: "success - display name with arabic",
			user: model.AuthUser{
				User: model.User{
					Login:       "arabicuser",
					DisplayName: "المستخدم أحمد",
				},
				PasswordHash: "hash_arabic",
			},
			expectedErr: nil,
		},
		{
			name: "success - display name with emoji",
			user: model.AuthUser{
				User: model.User{
					Login:       "emojiuser",
					DisplayName: "User 🚀 Test",
				},
				PasswordHash: "hash_emoji",
			},
			expectedErr: nil,
		},
		{
			name: "success - display name with mixed unicode",
			user: model.AuthUser{
				User: model.User{
					Login:       "mixeduser",
					DisplayName: "Mix Микс 混合 مختلط 😊",
				},
				PasswordHash: "hash_mixed",
			},
			expectedErr: nil,
		},
		{
			name: "error - duplicate login",
			user: model.AuthUser{
				User: model.User{
					Login:       "duplicateuser",
					DisplayName: "Duplicate User",
				},
				PasswordHash: "hash789",
			},
			expectedErr: derrors.ErrUserAlreadyExists,
		},
		{
			name: "success - empty login allowed by DB",
			user: model.AuthUser{
				User: model.User{
					Login:       "",
					DisplayName: "Empty Login User",
				},
				PasswordHash: "hash_empty_login",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			if errors.Is(tc.expectedErr, derrors.ErrUserAlreadyExists) {
				existingUser := model.AuthUser{
					User: model.User{
						Login:       tc.user.Login,
						DisplayName: "Existing User",
					},
					PasswordHash: "existing_hash",
				}
				_, err := s.repo.InsertUser(ctx, existingUser)
				s.Require().NoError(err)
			}

			// Test.
			result, err := s.repo.InsertUser(ctx, tc.user)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.NotZero(result.ID)
			s.Equal(tc.user.Login, result.Login)
			s.Equal(tc.user.PasswordHash, result.PasswordHash)
			s.Equal(tc.user.DisplayName, result.DisplayName)
			s.NotZero(result.CreatedAt)
			s.NotZero(result.UpdatedAt)
			s.True(result.DeletedAt.IsZero())
		})
	}
}

func (s *RepositoryTestSuite) TestSelectAuthUserByLogin() {
	testCases := []struct {
		expectedErr error
		name        string
		login       string
	}{
		{
			name:        "success - finds existing user",
			login:       "existinguser",
			expectedErr: nil,
		},
		{
			name:        "error - user not found",
			login:       "nonexistent",
			expectedErr: derrors.ErrUserNotFound,
		},
		{
			name:        "error - empty login returns not found",
			login:       "",
			expectedErr: derrors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var expectedUser model.AuthUser
			if tc.expectedErr == nil {
				expectedUser = model.AuthUser{
					User: model.User{
						Login:       tc.login,
						DisplayName: "Test User",
					},
					PasswordHash: "hash_" + tc.login,
				}
				insertedUser, err := s.repo.InsertUser(ctx, expectedUser)
				s.Require().NoError(err)

				expectedUser = insertedUser
			}

			// Test.
			result, err := s.repo.SelectAuthUserByLogin(ctx, tc.login)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.assertAuthUserEqual(expectedUser, result)
		})
	}
}

func (s *RepositoryTestSuite) TestSelectAuthUserByID() {
	testCases := []struct {
		expectedErr error
		name        string
		loginSuffix string
		createUser  bool
	}{
		{
			name:        "success - finds existing user",
			loginSuffix: "byid1",
			createUser:  true,
			expectedErr: nil,
		},
		{
			name:        "error - user not found",
			loginSuffix: "byid2",
			createUser:  false,
			expectedErr: derrors.ErrUserNotFound,
		},
		{
			name:        "error - zero ID returns not found",
			loginSuffix: "byid3",
			createUser:  false,
			expectedErr: derrors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var (
				expectedUser model.AuthUser
				queryUserID  int64
			)

			switch {
			case tc.createUser:
				expectedUser = model.AuthUser{
					User: model.User{
						Login:       "testuser_" + tc.loginSuffix,
						DisplayName: "Test User By ID",
					},
					PasswordHash: "hash_byid",
				}
				insertedUser, err := s.repo.InsertUser(ctx, expectedUser)
				s.Require().NoError(err)

				expectedUser = insertedUser
				queryUserID = insertedUser.ID
			case tc.name == "error - user not found":
				queryUserID = 999999
			default:
				queryUserID = 0
			}

			// Test.
			result, err := s.repo.SelectAuthUserByID(ctx, queryUserID)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.assertAuthUserEqual(expectedUser, result)
		})
	}
}

func (s *RepositoryTestSuite) TestSelectUserByID() {
	testCases := []struct {
		expectedErr error
		name        string
		loginSuffix string
		createUser  bool
	}{
		{
			name:        "success - finds existing user",
			loginSuffix: "seluser1",
			createUser:  true,
			expectedErr: nil,
		},
		{
			name:        "error - user not found",
			loginSuffix: "seluser2",
			createUser:  false,
			expectedErr: derrors.ErrUserNotFound,
		},
		{
			name:        "error - zero ID returns not found",
			loginSuffix: "seluser3",
			createUser:  false,
			expectedErr: derrors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var (
				expectedUser model.User
				queryUserID  int64
			)

			switch {
			case tc.createUser:
				authUser := model.AuthUser{
					User: model.User{
						Login:       "selectuser_" + tc.loginSuffix,
						DisplayName: "Select User",
					},
					PasswordHash: "hash_select",
				}
				insertedUser, err := s.repo.InsertUser(ctx, authUser)
				s.Require().NoError(err)

				expectedUser = insertedUser.User
				queryUserID = insertedUser.ID
			case tc.name == "error - user not found":
				queryUserID = 999999
			default:
				queryUserID = 0
			}

			// Test.
			result, err := s.repo.SelectUserByID(ctx, queryUserID)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.assertUserEqual(expectedUser, result)
		})
	}
}

func (s *RepositoryTestSuite) TestUpdateUser() {
	testCases := []struct {
		expectedErr  error
		name         string
		login        string
		expectedName string
		expectedHash string
		updateUser   model.AuthUser
	}{
		{
			name:  "success - updates display name and password hash",
			login: "updateuser1",
			updateUser: model.AuthUser{
				User: model.User{
					DisplayName: "Updated Name",
				},
				PasswordHash: "new_hash",
			},
			expectedErr:  nil,
			expectedName: "Updated Name",
			expectedHash: "new_hash",
		},
		{
			name:  "success - updates to empty display name",
			login: "updateuser2",
			updateUser: model.AuthUser{
				User: model.User{
					DisplayName: "",
				},
				PasswordHash: "hash_empty_name",
			},
			expectedErr:  nil,
			expectedName: "",
			expectedHash: "hash_empty_name",
		},
		{
			name:  "error - user not found",
			login: "updateuser3",
			updateUser: model.AuthUser{
				User: model.User{
					ID:          999999,
					DisplayName: "Not Found",
				},
				PasswordHash: "hash_notfound",
			},
			expectedErr: derrors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var originalUser, updateUser model.AuthUser
			if tc.expectedErr == nil {
				originalUser = model.AuthUser{
					User: model.User{
						Login:       tc.login,
						DisplayName: "Original Name",
					},
					PasswordHash: "original_hash",
				}
				insertedUser, err := s.repo.InsertUser(ctx, originalUser)
				s.Require().NoError(err)

				updateUser = tc.updateUser
				updateUser.ID = insertedUser.ID
				updateUser.Login = insertedUser.Login
				updateUser.CreatedAt = insertedUser.CreatedAt
			}

			// Test.
			result, err := s.repo.UpdateUser(ctx, updateUser)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(result.UpdatedAt)

				return
			}

			s.Require().NoError(err)
			s.Equal(updateUser.ID, result.ID)
			s.Equal(tc.expectedName, result.DisplayName)
			s.Equal(tc.expectedHash, result.PasswordHash)
			s.NotZero(result.UpdatedAt)
			s.True(result.UpdatedAt.After(originalUser.UpdatedAt))
		})
	}
}

func (s *RepositoryTestSuite) TestDeleteUser() {
	testCases := []struct {
		expectedErr error
		name        string
		login       string
		userID      int64
	}{
		{
			name:        "success - deletes existing user",
			login:       "deleteuser1",
			userID:      0,
			expectedErr: nil,
		},
		{
			name:        "error - user not found",
			login:       "deleteuser2",
			userID:      999999,
			expectedErr: derrors.ErrUserNotFound,
		},
		{
			name:        "error - already deleted user",
			login:       "deleteuser3",
			userID:      0,
			expectedErr: derrors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			userID := tc.userID
			if tc.login != "" && userID == 0 {
				user := model.AuthUser{
					User: model.User{
						Login:       tc.login,
						DisplayName: "User To Delete",
					},
					PasswordHash: "hash_delete",
				}
				insertedUser, err := s.repo.InsertUser(ctx, user)
				s.Require().NoError(err)

				userID = insertedUser.ID

				if errors.Is(tc.expectedErr, derrors.ErrUserNotFound) && tc.name == "error - already deleted user" {
					err = s.repo.DeleteUser(ctx, userID)
					s.Require().NoError(err)
				}
			}

			// Test.
			err := s.repo.DeleteUser(ctx, userID)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				return
			}

			s.Require().NoError(err)

			_, err = s.repo.SelectAuthUserByID(ctx, userID)
			s.ErrorIs(err, derrors.ErrUserNotFound)
		})
	}
}

func (s *RepositoryTestSuite) TestSelectUsersByLoginPrefix() {
	testCases := []struct {
		name           string
		excludeLogin   string
		prefix         string
		setupUsers     []string
		expectedLogins []string
		limit          int
		offset         int
	}{
		{
			name:           "success - finds users by prefix",
			setupUsers:     []string{"search_alice", "search_bob", "search_charlie", "other_user"},
			excludeLogin:   "other_user",
			prefix:         "search_",
			limit:          10,
			offset:         0,
			expectedLogins: []string{"search_alice", "search_bob", "search_charlie"},
		},
		{
			name:           "success - pagination with limit",
			setupUsers:     []string{"page_user", "page_a", "page_b", "page_c", "page_d"},
			excludeLogin:   "page_user",
			prefix:         "page_",
			limit:          2,
			offset:         0,
			expectedLogins: []string{"page_a", "page_b"},
		},
		{
			name:           "success - pagination with offset",
			setupUsers:     []string{"off_user", "off_a", "off_b", "off_c", "off_d"},
			excludeLogin:   "off_user",
			prefix:         "off_",
			limit:          2,
			offset:         2,
			expectedLogins: []string{"off_c", "off_d"},
		},
		{
			name:           "success - empty result when no matches",
			setupUsers:     []string{"nomatch_user", "nomatch_a", "nomatch_b"},
			excludeLogin:   "nomatch_user",
			prefix:         "different_",
			limit:          10,
			offset:         0,
			expectedLogins: []string{},
		},
		{
			name:           "success - excludes current user",
			setupUsers:     []string{"exclude_me", "exclude_other"},
			excludeLogin:   "exclude_me",
			prefix:         "exclude_",
			limit:          10,
			offset:         0,
			expectedLogins: []string{"exclude_other"},
		},
		{
			name:           "success - prefix matches multiple words",
			setupUsers:     []string{"multi_current", "multi_alice", "multi_bob"},
			excludeLogin:   "multi_current",
			prefix:         "multi_",
			limit:          10,
			offset:         0,
			expectedLogins: []string{"multi_alice", "multi_bob"},
		},
		{
			name:           "success - case sensitive prefix",
			setupUsers:     []string{"case_user", "case_lower", "CASE_UPPER"},
			excludeLogin:   "case_user",
			prefix:         "case_",
			limit:          10,
			offset:         0,
			expectedLogins: []string{"case_lower"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var currentUserID int64

			for _, login := range tc.setupUsers {
				user := model.AuthUser{
					User: model.User{
						Login:       login,
						DisplayName: "User " + login,
					},
					PasswordHash: "hash_" + login,
				}
				insertedUser, err := s.repo.InsertUser(ctx, user)
				s.Require().NoError(err)

				if login == tc.excludeLogin {
					currentUserID = insertedUser.ID
				}
			}

			// Test.
			result, err := s.repo.SelectUsersByLoginPrefix(ctx, currentUserID, tc.prefix, tc.limit, tc.offset)

			// Check.
			s.Require().NoError(err)
			s.Len(result, len(tc.expectedLogins))

			for i, expectedLogin := range tc.expectedLogins {
				s.Equal(expectedLogin, result[i].Login)
				s.NotZero(result[i].ID)
				s.NotEmpty(result[i].DisplayName)
				s.NotZero(result[i].CreatedAt)
				s.NotEqual(currentUserID, result[i].ID)
			}

			if len(result) > 1 {
				for i := 1; i < len(result); i++ {
					s.Less(result[i-1].Login, result[i].Login)
				}
			}
		})
	}
}

func (s *RepositoryTestSuite) TestCountUsersByLoginPrefix() {
	testCases := []struct {
		name          string
		excludeLogin  string
		prefix        string
		setupUsers    []string
		expectedCount int
	}{
		{
			name:          "success - counts users by prefix",
			setupUsers:    []string{"count_alice", "count_bob", "count_charlie", "other_user"},
			excludeLogin:  "other_user",
			prefix:        "count_",
			expectedCount: 3,
		},
		{
			name:          "success - zero count when no matches",
			setupUsers:    []string{"zero_user", "zero_a"},
			excludeLogin:  "zero_user",
			prefix:        "different_",
			expectedCount: 0,
		},
		{
			name:          "success - excludes current user from count",
			setupUsers:    []string{"excl_me", "excl_other"},
			excludeLogin:  "excl_me",
			prefix:        "excl_",
			expectedCount: 1,
		},
		{
			name:          "success - counts with multiple matching users",
			setupUsers:    []string{"many_current", "many_a", "many_b", "many_c", "many_d", "many_e"},
			excludeLogin:  "many_current",
			prefix:        "many_",
			expectedCount: 5,
		},
		{
			name:          "success - case sensitive count",
			setupUsers:    []string{"sens_user", "sens_lower", "SENS_UPPER"},
			excludeLogin:  "sens_user",
			prefix:        "sens_",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			var currentUserID int64

			for _, login := range tc.setupUsers {
				user := model.AuthUser{
					User: model.User{
						Login:       login,
						DisplayName: "User " + login,
					},
					PasswordHash: "hash_" + login,
				}
				insertedUser, err := s.repo.InsertUser(ctx, user)
				s.Require().NoError(err)

				if login == tc.excludeLogin {
					currentUserID = insertedUser.ID
				}
			}

			// Test.
			count, err := s.repo.CountUsersByLoginPrefix(ctx, currentUserID, tc.prefix)

			// Check.
			s.Require().NoError(err)
			s.Equal(tc.expectedCount, count)
		})
	}
}
