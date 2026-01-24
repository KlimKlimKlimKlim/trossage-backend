package hasher

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
)

type HasherTestSuite struct {
	suite.Suite
}

func TestHasherTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HasherTestSuite))
}

func (s *HasherTestSuite) TestNew() {
	testCases := []struct {
		config         *config.Hasher
		name           string
		expectedMemory uint32
		expectedIter   uint32
		expectedSalt   uint32
		expectedKey    uint32
	}{
		{
			name: "success - creates hasher with default config",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			expectedMemory: 65536,
			expectedIter:   3,
			expectedSalt:   16,
			expectedKey:    32,
		},
		{
			name: "success - creates hasher with custom config",
			config: &config.Hasher{
				Memory:     131072,
				Iterations: 5,
				SaltLength: 32,
				KeyLength:  64,
			},
			expectedMemory: 131072,
			expectedIter:   5,
			expectedSalt:   32,
			expectedKey:    64,
		},
		{
			name: "success - creates hasher with minimum values",
			config: &config.Hasher{
				Memory:     8192,
				Iterations: 1,
				SaltLength: 8,
				KeyLength:  16,
			},
			expectedMemory: 8192,
			expectedIter:   1,
			expectedSalt:   8,
			expectedKey:    16,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.T().Parallel()

			// Test.
			hasher := New(tc.config)

			// Check.
			s.Require().NotNil(hasher)
			s.Require().NotNil(hasher.params)
			s.Equal(tc.expectedMemory, hasher.params.Memory)
			s.Equal(tc.expectedIter, hasher.params.Iterations)
			s.Equal(tc.expectedSalt, hasher.params.SaltLength)
			s.Equal(tc.expectedKey, hasher.params.KeyLength)

			cpuCount := runtime.NumCPU()

			var expectedParallelism uint8

			switch {
			case cpuCount > maxUint8:
				expectedParallelism = maxUint8
			case cpuCount < 1:
				expectedParallelism = 1
			default:
				expectedParallelism = uint8(cpuCount)
			}

			s.Equal(expectedParallelism, hasher.params.Parallelism)
			s.GreaterOrEqual(hasher.params.Parallelism, uint8(1))
			s.LessOrEqual(hasher.params.Parallelism, uint8(maxUint8))
		})
	}
}

//nolint:gosmopolitan // Testing Unicode support in password hashing
func (s *HasherTestSuite) TestHashPassword() {
	testCases := []struct {
		name              string
		config            *config.Hasher
		password          string
		testDifferentSalt bool
	}{
		{
			name: "success - hashes simple password",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "password123",
		},
		{
			name: "success - hashes empty password",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "",
		},
		{
			name: "success - hashes password with special characters",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "P@ssw0rd!#$%^&*()",
		},
		{
			name: "success - hashes password with spaces",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "password with spaces",
		},
		{
			name: "success - hashes unicode password cyrillic",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "Пароль123",
		},
		{
			name: "success - hashes unicode password chinese",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "密码123",
		},
		{
			name: "success - hashes unicode password arabic",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "كلمةالسر123",
		},
		{
			name: "success - hashes unicode password emoji",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "pass🔒word🔑123",
		},
		{
			name: "success - hashes very long password",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: strings.Repeat("a", 1000),
		},
		{
			name: "success - hashes password with newlines",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "password\nwith\nnewlines",
		},
		{
			name: "success - hashes password with tabs",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password: "password\twith\ttabs",
		},
		{
			name: "success - low security config",
			config: &config.Hasher{
				Memory:     8192,
				Iterations: 1,
				SaltLength: 8,
				KeyLength:  16,
			},
			password: "testpassword",
		},
		{
			name: "success - high security config",
			config: &config.Hasher{
				Memory:     262144,
				Iterations: 5,
				SaltLength: 32,
				KeyLength:  64,
			},
			password: "testpassword",
		},
		{
			name: "success - same password produces different salts",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:          "testpassword",
			testDifferentSalt: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.T().Parallel()

			hasher := New(tc.config)

			// Test.
			hash, err := hasher.HashPassword(tc.password)

			// Check.
			s.Require().NoError(err)
			s.NotEmpty(hash)
			s.Contains(hash, "$argon2id$")
			s.Greater(len(hash), 50)

			if tc.testDifferentSalt {
				hash2, err := hasher.HashPassword(tc.password)
				s.Require().NoError(err)

				hash3, err := hasher.HashPassword(tc.password)
				s.Require().NoError(err)

				s.NotEqual(hash, hash2)
				s.NotEqual(hash2, hash3)
				s.NotEqual(hash, hash3)
			}
		})
	}
}

//nolint:gosmopolitan // Testing Unicode support in password verification
func (s *HasherTestSuite) TestVerifyPassword() {
	testCases := []struct {
		name              string
		config            *config.Hasher
		password          string
		wrongPassword     string
		invalidHash       string
		testCaseSensitive bool
		testWhitespace    bool
		expectError       bool
	}{
		{
			name: "success - verifies simple password",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:      "password123",
			wrongPassword: "wrongpassword",
		},
		{
			name: "success - verifies empty password",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:      "",
			wrongPassword: "notempty",
		},
		{
			name: "success - verifies password with special characters",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:      "P@ssw0rd!#$%",
			wrongPassword: "P@ssw0rd!#$",
		},
		{
			name: "success - verifies unicode password cyrillic",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:      "Пароль123",
			wrongPassword: "Пароль124",
		},
		{
			name: "success - verifies unicode password chinese",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:      "密码123",
			wrongPassword: "密码124",
		},
		{
			name: "success - verifies unicode password emoji",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:      "pass🔒word🔑",
			wrongPassword: "pass🔒word🗝",
		},
		{
			name: "success - verifies long password",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:      strings.Repeat("a", 500),
			wrongPassword: strings.Repeat("b", 500),
		},
		{
			name: "error - completely invalid hash",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:    "password123",
			invalidHash: "notahash",
			expectError: true,
		},
		{
			name: "error - empty hash",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:    "password123",
			invalidHash: "",
			expectError: true,
		},
		{
			name: "error - partial argon2id format",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:    "password123",
			invalidHash: "$argon2id$",
			expectError: true,
		},
		{
			name: "error - wrong algorithm format",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:    "password123",
			invalidHash: "$bcrypt$12$abc123def456",
			expectError: true,
		},
		{
			name: "error - malformed argon2id hash",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:    "password123",
			invalidHash: "$argon2id$v=19$m=65536$invalid",
			expectError: true,
		},
		{
			name: "success - password verification is case sensitive",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:          "Password123",
			testCaseSensitive: true,
		},
		{
			name: "success - password verification treats whitespace strictly",
			config: &config.Hasher{
				Memory:     65536,
				Iterations: 3,
				SaltLength: 16,
				KeyLength:  32,
			},
			password:       "password 123",
			testWhitespace: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.T().Parallel()

			hasher := New(tc.config)

			if tc.invalidHash != "" {
				// Test.
				match, err := hasher.VerifyPassword(tc.password, tc.invalidHash)

				// Check.
				if tc.expectError {
					s.Require().Error(err)
					s.False(match)
				} else {
					s.Require().NoError(err)
				}

				return
			}

			hash, err := hasher.HashPassword(tc.password)
			s.Require().NoError(err)

			// Test - correct password.
			match, err := hasher.VerifyPassword(tc.password, hash)

			// Check.
			s.Require().NoError(err)
			s.True(match)

			if tc.testCaseSensitive {
				// Test - lowercase.
				match, err = hasher.VerifyPassword("password123", hash)
				s.Require().NoError(err)
				s.False(match)

				// Test - uppercase.
				match, err = hasher.VerifyPassword("PASSWORD123", hash)
				s.Require().NoError(err)
				s.False(match)

				return
			}

			if tc.testWhitespace {
				// Test - no space.
				match, err = hasher.VerifyPassword("password123", hash)
				s.Require().NoError(err)
				s.False(match)

				// Test - extra space.
				match, err = hasher.VerifyPassword("password  123", hash)
				s.Require().NoError(err)
				s.False(match)

				// Test - trailing space.
				match, err = hasher.VerifyPassword("password 123 ", hash)
				s.Require().NoError(err)
				s.False(match)

				return
			}

			if tc.wrongPassword != "" {
				// Test - wrong password.
				match, err = hasher.VerifyPassword(tc.wrongPassword, hash)

				// Check.
				s.Require().NoError(err)
				s.False(match)
			}
		})
	}
}

func (s *HasherTestSuite) TestCrossConfigVerification() {
	s.T().Parallel()

	config1 := &config.Hasher{
		Memory:     65536,
		Iterations: 3,
		SaltLength: 16,
		KeyLength:  32,
	}

	config2 := &config.Hasher{
		Memory:     131072,
		Iterations: 5,
		SaltLength: 32,
		KeyLength:  64,
	}

	hasher1 := New(config1)
	hasher2 := New(config2)

	password := "testpassword123"

	hash, err := hasher1.HashPassword(password)
	s.Require().NoError(err)

	match, err := hasher1.VerifyPassword(password, hash)
	s.Require().NoError(err)
	s.True(match, "Same hasher should verify its own hash")

	match, err = hasher2.VerifyPassword(password, hash)
	s.Require().NoError(err)
	s.True(match, "Different config hasher should still verify the hash because params are in the hash")

	match, err = hasher2.VerifyPassword("wrongpassword", hash)
	s.Require().NoError(err)
	s.False(match, "Wrong password should not match")
}
