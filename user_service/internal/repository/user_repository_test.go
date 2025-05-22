package repository

import (
	"context"
	"errors"
	"testing"
	"user_service/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestCreateUser_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	user := &domain.User{
		ID:       uuid.New().String(),
		Username: "testuser",
		Password: "testpassword",
		Name:    "Test User",
		Email:  "testuser@gmail.com",
		Role:  "user",
	}

	// Mock the database query
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users" \("id","username","password","name","email","role"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
		WithArgs(user.ID, user.Username, user.Password, user.Name, user.Email, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	userRepo := NewUserRepository(db)

	userID, err := userRepo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotEmpty(t, userID)
	assert.Equal(t, user.ID, userID)
}

func TestCreateUser_Failure(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	user := &domain.User{
		ID:       uuid.New().String(),
		Username: "testuser",
		Password: "testpassword",
		Name:    "Test User",
		Email:  "",
		Role:  "user",
	}

	// Mock the database query to return an error
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users" \("id","username","password","name","email","role"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
		WithArgs(user.ID, user.Username, user.Password, user.Name, user.Email, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	userRepo := NewUserRepository(db)

	userID, err := userRepo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotEmpty(t, userID)


	user2 := &domain.User{
		ID:       "22222222-2222-2222-2222-222222222222",
		Username: "testuser1", // duplicate username
		Password: "pass2",
		Name:     "User Two",
		Email:    "", // duplicate email
		Role:     "admin",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users" \("id","username","password","name","email","role"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
		WithArgs(user2.ID, user2.Username, user2.Password, user2.Name, user2.Email, user2.Role).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))
	mock.ExpectRollback()

	id2, err2 := userRepo.CreateUser(context.Background(), user2)
	assert.Error(t, err2)
	assert.Empty(t, id2)
	assert.Contains(t, err2.Error(), "duplicate key")
}

func TestLogin_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	user := &domain.User{
		ID:       uuid.New().String(),
		Username: "testuser",
		Password: "testpassword",
		Name:    "Test User",
		Email:  "testuser@gmail.com",
		Role:  "user",
	}

	// Mock the database query
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(user.Username, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"}).
		AddRow(user.ID, user.Username, user.Password, user.Name, user.Email, user.Role))
	
	userRepo := NewUserRepository(db)

	// Call the Login method
	loggedInUser, err := userRepo.Login(context.Background(), user.Username)
	assert.NoError(t, err)
	assert.NotNil(t, loggedInUser)
	assert.Equal(t, user.ID, loggedInUser.ID)
	assert.Equal(t, user.Username, loggedInUser.Username)
	assert.Equal(t, user.Password, loggedInUser.Password)
	assert.Equal(t, user.Name, loggedInUser.Name)
	assert.Equal(t, user.Email, loggedInUser.Email)
	assert.Equal(t, user.Role, loggedInUser.Role)
}

func TestLogin_Failure(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	// Mock the database query to return no rows
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"."id" LIMIT \$2`).
		WithArgs("nonexistentuser", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"}))
	
	userRepo := NewUserRepository(db)

	// Call the Login method with a non-existent username
	loggedInUser, err := userRepo.Login(context.Background(), "nonexistentuser")
	assert.Error(t, err)
	assert.Nil(t, loggedInUser)
	assert.Equal(t, "record not found", err.Error())
}

func TestGetUserByID_Success(t *testing.T){
	db, mock := setupMockDB(t)
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	expectedUser := &domain.User{
		ID:       uuid.New().String(),
		Username: "testuser",
		Password: "testpassword",
		Name:    "Test User",
		Email:  "testemail@gmail.com",
		Role:  "user",
	}

	// Mock the database query
	// insert the expected user into the mock database
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(expectedUser.ID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password, expectedUser.Name, expectedUser.Email, expectedUser.Role))

	userRepo := NewUserRepository(db)

	// Call the GetUserByID method
	user, err := userRepo.GetUserByID(context.Background(), expectedUser.ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Role, user.Role)
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	// Mock the database query to return no rows
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"."id" LIMIT \$2`).
		WithArgs("nonexistent-id", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"}))

	userRepo := NewUserRepository(db)

	// Call the GetUserByID method with a non-existent ID
	user, err := userRepo.GetUserByID(context.Background(), "nonexistent-id")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "record not found", err.Error())
}

func TestGetAllUsers_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	expectedUsers := []*domain.User{
		{
			ID:       uuid.New().String(),
			Username: "testuser1",
			Password: "testpassword1",
			Name:    "Test User 1",
			Email:  "testuser1@gmail.com",
			Role:  "user",
		},
		{
			ID:       uuid.New().String(),
			Username: "testuser2",
			Password: "testpassword2",
			Name:    "Test User 2",
			Email:  "testuser2@gmail.com",
			Role:  "admin",
		},
	}

	// Mock the database query
	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"}).
		AddRow(expectedUsers[0].ID, expectedUsers[0].Username, expectedUsers[0].Password, expectedUsers[0].Name, expectedUsers[0].Email, expectedUsers[0].Role).
		AddRow(expectedUsers[1].ID, expectedUsers[1].Username, expectedUsers[1].Password, expectedUsers[1].Name, expectedUsers[1].Email, expectedUsers[1].Role))
	
	userRepo := NewUserRepository(db)

	// Call the GetAllUsers method
	users, err := userRepo.GetAllUsers(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	
	for i, user := range users {
		assert.Equal(t, expectedUsers[i].ID, user.ID)
		assert.Equal(t, expectedUsers[i].Username, user.Username)
		assert.Equal(t, expectedUsers[i].Password, user.Password)
		assert.Equal(t, expectedUsers[i].Name, user.Name)
		assert.Equal(t, expectedUsers[i].Email, user.Email)
		assert.Equal(t, expectedUsers[i].Role, user.Role)
	}
}

func TestGetAllUsers_Failure(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	// Mock the database query to return an error
	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnError(errors.New("database error"))

	userRepo := NewUserRepository(db)

	// Call the GetAllUsers method
	users, err := userRepo.GetAllUsers(context.Background())
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, "database error", err.Error())
}