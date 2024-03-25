package postgres

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/taldoflemis/brain.test/test/helpers"
)

type LocalIDPPostgresStorerTestSuite struct {
	suite.Suite
	pgContainer *testshelpers.PostgresContainer
	ctx         context.Context
	repo        *LocalIDPPostgresStorer
	pool        *pgxpool.Pool
}

func (suite *LocalIDPPostgresStorerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testshelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	pool, err := NewPool(pgContainer.ConnStr)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(pgContainer.ConnStr, "./migrations/")

	repository := NewLocalIDPPostgresStorer(pool)

	suite.repo = repository
	suite.pool = pool
}

func (suite *LocalIDPPostgresStorerTestSuite) SetupTest() {
	_, err := suite.pool.Exec(
		suite.ctx,
		`INSERT INTO users (id, username, email, password) VALUES ('f7396104-a636-4826-9d9f-b92ae90cea14', 'gepeto', 'gepeto@gmail.com', 'hashedpass')`,
	)
	if err != nil {
		log.Fatalf("error inserting user: %s", err)
	}
}

func (suite *LocalIDPPostgresStorerTestSuite) TearDownTest() {
	_, err := suite.pool.Exec(suite.ctx, "TRUNCATE TABLE users")
	if err != nil {
		log.Fatalf("error truncating users table: %s", err)
	}
}

func (suite *LocalIDPPostgresStorerTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestLocalIDPPostgresStorer(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	suite.Run(t, new(LocalIDPPostgresStorerTestSuite))
}

func (suite *LocalIDPPostgresStorerTestSuite) TestFindUserById() {
	// Arrange
	t := suite.T()

	// Act
	user, err := suite.repo.FindUserByUsername(suite.ctx, "gepeto")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "gepeto", user.Username)
	assert.Equal(t, "gepeto@gmail.com", user.Email)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestTryToFindUserByIdThatDoesNotExist() {
	// Arrange
	t := suite.T()
	random := "7ebc4755-b7cc-4963-a2b1-636949b035d6"

	// Act
	user, err := suite.repo.FindUserByUsername(suite.ctx, random)

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestDeleteUser() {
	// Arrange
	t := suite.T()

	// Act
	err := suite.repo.DeleteUser(suite.ctx, "f7396104-a636-4826-9d9f-b92ae90cea14")

	// Assert
	assert.NoError(t, err)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestTryToDeleteUserThatDoesNotExist() {
	// Arrange
	t := suite.T()
	randomId := "d0b8b515-f46b-4179-bb26-f7833ded8f8f"

	// Act
	err := suite.repo.DeleteUser(suite.ctx, randomId)

	// Assert
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestCreateUser() {
	// Arrange
	t := suite.T()
	username := "tubias"
	email := "tubias2@gmail.com"
	password := "hashedpassword"

	// Act
	user, err := suite.repo.StoreUser(suite.ctx, username, email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestUpdateUser() {
	// Arrange
	t := suite.T()
	username := "tubias"
	email := "tubias3@gmail.com"
	password := "hashedpassword"
	userId := "f7396104-a636-4826-9d9f-b92ae90cea14"

	// Act
	user, err := suite.repo.UpdateUser(suite.ctx, userId, username, password, email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestTryToUpdateUserThatDoesNotExist() {
	// Arrange
	t := suite.T()
	username := "tubias"
	email := "tubias@gmail.com"
	password := "hashedpassword"
	randomId := "d0b8b515-f46b-4179-bb26-f7833ded8f8f"

	// Act
	user, err := suite.repo.UpdateUser(suite.ctx, randomId, username, password, email)

	// Assert
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, user)
}
