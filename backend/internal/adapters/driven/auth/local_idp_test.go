package auth

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
	testshelpers "github.com/taldoflemis/brain.test/test/helpers"
)

const (
	seed = "SzceVsT4GdFOlrZn60XMgrFcvMNUMuuJ"
)

var (
	testUsername         = "gepeto"
	testEmail            = "gepeto@gmail.com"
	testPassword         = "mypassword"
	testHashedPassword   = "$2a$12$TSjLw2cqeD5bcjPUgOWaaew3xP88soPytNTnMi27vxcNMCDaLFkBa"
	testUserId           = ""
	accessMaxAgeInMin    = 15
	refreshMaxAgeInHours = 24
)

type LocalIDPTestSuite struct {
	suite.Suite
	pgContainer *testshelpers.PostgresContainer
	ctx         context.Context
	svc         *localIDP
	repo        *postgres.LocalIDPPostgresStorer
	pool        *pgxpool.Pool
}

func (suite *LocalIDPTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testshelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	pool, err := postgres.NewPool(pgContainer.ConnStr)
	if err != nil {
		log.Fatal(err)
	}

	postgres.Migrate(pgContainer.ConnStr, "../postgres/migrations/")

	repository := postgres.NewLocalIDPPostgresStorer(pool)
	logger := testshelpers.NewDummyLogger(log.Writer())
	cfg := NewLocalIdpConfig(seed, "issuer", "audience", accessMaxAgeInMin, refreshMaxAgeInHours)
	svc := NewLocalIdp(*cfg, logger, repository)

	suite.svc = svc
	suite.repo = repository
	suite.pool = pool
}

func (suite *LocalIDPTestSuite) SetupTest() {
	user, err := suite.repo.StoreUser(suite.ctx, testUsername, testEmail, testHashedPassword)
	if err != nil {
		log.Fatalf("error inserting user: %s", err)
	}
	testUserId = user.ID
}

func (suite *LocalIDPTestSuite) TearDownTest() {
	_, err := suite.pool.Exec(suite.ctx, "TRUNCATE TABLE users")
	if err != nil {
		log.Fatalf("error truncating users table: %s", err)
	}
}

func (suite *LocalIDPTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestLocalIDP(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	suite.Run(t, new(LocalIDPTestSuite))
}

func (suite *LocalIDPTestSuite) TestCreateUser() {
	// Arrange
	t := suite.T()
	username := "tubias"
	email := "tubias@gmail.com"
	password := "hashedpass"

	// Act
	user, err := suite.svc.CreateUser(suite.ctx, username, email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
}

func (suite *LocalIDPTestSuite) TestCreateToken() {
	// Arrange
	t := suite.T()

	// Act
	tokenResponse, err := suite.svc.CreateToken(suite.ctx, testUserId)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenResponse.AccessToken)
	assert.NotEmpty(t, tokenResponse.RefreshToken)
	assert.WithinDurationf(
		t,
		time.Now().Add((time.Duration(accessMaxAgeInMin) * time.Minute)),
		tokenResponse.ExpiresAt,
		5*time.Second,
		"expected expiresAt to be within %d minutes",
		accessMaxAgeInMin,
	)
}

func (suite *LocalIDPTestSuite) TestRefreshToken() {
	// Arrange
	t := suite.T()
	validRefreshToken, err := suite.svc.generateToken(suite.ctx, testUserId, time.Hour)
	assert.NoError(t, err)

	// Act
	tokenResponse, err := suite.svc.RefreshToken(suite.ctx, validRefreshToken)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenResponse.AccessToken)
	assert.NotEmpty(t, tokenResponse.RefreshToken)
	assert.WithinDurationf(
		t,
		time.Now().Add((time.Duration(accessMaxAgeInMin) * time.Minute)),
		tokenResponse.ExpiresAt,
		5*time.Second,
		"expected expiresAt to be within %d minutes",
		accessMaxAgeInMin,
	)
}

func (suite *LocalIDPTestSuite) TestInvalidRefreshToken() {
	// Arrange
	t := suite.T()
	invalidRefreshToken := "invalid"

	// Act
	tokenResponse, err := suite.svc.RefreshToken(suite.ctx, invalidRefreshToken)

	// Assert
	assert.ErrorIs(t, err, ErrInvalidRefreshToken)
	assert.Nil(t, tokenResponse)
}

func (suite *LocalIDPTestSuite) TestExpiredRefreshToken() {
	// Arrange
	t := suite.T()
	expiredRefreshToken, err := suite.svc.generateToken(suite.ctx, testUserId, -time.Minute)
	assert.NoError(t, err)

	// Act
	tokenResponse, err := suite.svc.RefreshToken(suite.ctx, expiredRefreshToken)

	// Assert
	assert.ErrorIs(t, err, ErrExpiredToken)
	assert.Nil(t, tokenResponse)
}

func (suite *LocalIDPTestSuite) TestAuthenticateUser() {
	// Arrange
	t := suite.T()

	// Act
	user, err := suite.svc.AuthenticateUser(suite.ctx, testUsername, testPassword)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUsername, user.Username)
	assert.Equal(t, testEmail, user.Email)
}

func (suite *LocalIDPTestSuite) TestAuthenticateUserWithWrongPassword() {
	// Arrange
	t := suite.T()
	wrongPassword := "wrongpassword"

	// Act
	user, err := suite.svc.AuthenticateUser(suite.ctx, testUsername, wrongPassword)

	// Assert
	assert.ErrorIs(t, err, ErrInvalidPassword)
	assert.Nil(t, user)
}

func (suite *LocalIDPTestSuite) TestDeleteUser() {
	// Arrange
	t := suite.T()

	// Act
	err := suite.svc.DeleteUser(suite.ctx, testUserId)

	// Assert
	assert.NoError(t, err)
}

func (suite *LocalIDPTestSuite) TestDeleteUserThatDoesNotExist() {
	// Arrange
	t := suite.T()
	randomId := "d0b8b515-f46b-4179-bb26-f7833ded8f8f"

	// Act
	err := suite.svc.DeleteUser(suite.ctx, randomId)

	// Assert
	assert.ErrorContains(t, err, ErrUserNotFound.Error())
}

func (suite *LocalIDPTestSuite) TestUpdateUser() {
	// Arrange
	t := suite.T()
	username := "tubias"
	email := "tubias@gmail.com"
	password := "hashedpassword"

	// Act
	user, err := suite.svc.UpdateUser(suite.ctx, testUserId, username, password, email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
}

func (suite *LocalIDPTestSuite) TestUpdateUserThatDoesNotExist() {
	// Arrange
	t := suite.T()
	username := "tubias"
	email := "tubias@gmail.com"
	password := "hashedpassword"
	randomId := "d0b8b515-f46b-4179-bb26-f7833ded8f8f"

	// Act
	user, err := suite.svc.UpdateUser(suite.ctx, randomId, username, password, email)

	// Assert
	assert.ErrorContains(t, err, ErrUserNotFound.Error())
	assert.Nil(t, user)
}
