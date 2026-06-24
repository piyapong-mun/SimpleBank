package db

import (
	"context"
	"testing"
	"time"

	"github.com/piyapong-mun/simplebank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomString(12),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(15),
		Email:          util.RandomString(10) + "@email.com",
	}

	user, err := TestQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)

	assert.True(t, user.PasswordChangeAt.IsZero())
	assert.NotZero(t, user.CreateAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := TestQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.Email, user2.Email)
	assert.WithinDuration(t, user1.PasswordChangeAt, user2.PasswordChangeAt, time.Second)
	assert.WithinDuration(t, user1.CreateAt, user2.CreateAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		Username:         user1.Username,
		HashedPassword:   util.RandomString(20),
		FullName:         util.RandomString(15),
		Email:            util.RandomString(10) + "@new-email.com",
		PasswordChangeAt: time.Now().UTC(),
	}

	user2, err := TestQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, arg.HashedPassword, user2.HashedPassword)
	assert.Equal(t, arg.FullName, user2.FullName)
	assert.Equal(t, arg.Email, user2.Email)
	assert.WithinDuration(t, arg.PasswordChangeAt, user2.PasswordChangeAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := TestQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)

	user2, err := TestQueries.GetUser(context.Background(), user1.Username)
	require.Error(t, err)
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := TestQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	assert.True(t, len(users) >= 5)

	for _, user := range users {
		assert.NotEmpty(t, user)
	}
}
