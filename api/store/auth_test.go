package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

var (
	user1 = User{
		Login:        "ivanov",
		FirstName:    "Иван",
		LastName:     "Иванов",
		PasswordHash: "incorrect",
	}
	user2 = User{
		Login:        "petrov",
		FirstName:    "Петр",
		LastName:     "Петров",
		PasswordHash: "incorrect",
	}
	user3 = User{
		Login:        "sidorov",
		FirstName:    "Михаил",
		LastName:     "Сидоров",
		PasswordHash: "incorrect",
	}
)

func checkEqual(t *testing.T, first, second User) {
	assert.Equal(t, first.Login, second.Login)
	assert.Equal(t, first.FirstName, second.FirstName)
	assert.Equal(t, first.LastName, second.LastName)
	assert.Equal(t, first.PasswordHash, second.PasswordHash)
}

func TestCreateAndFindUser(t *testing.T) {
	initTestDatabase()
	defer removeTestDatabase()

	err := Auth.CreateUser(user1)
	assert.NoError(t, err)
	f1, err := Auth.FindUserByLogin(user1.Login)
	assert.NoError(t, err)
	checkEqual(t, user1, f1)
	err = Auth.CreateUser(user2)
	assert.NoError(t, err)
	err = Auth.CreateUser(user3)
	assert.NoError(t, err)
	f2, err := Auth.FindUserByLogin(user2.Login)
	assert.NoError(t, err)
	checkEqual(t, user2, f2)
	f3, err := Auth.FindUserByLogin(user3.Login)
	assert.NoError(t, err)
	checkEqual(t, user3, f3)
}

func TestUniqueLogin(t *testing.T) {
	initTestDatabase()
	defer removeTestDatabase()

	err := Auth.CreateUser(user1)
	assert.NoError(t, err)
	u2 := user2
	u2.Login = user1.Login
	err = Auth.CreateUser(user1)
	t.Log(xerrors.Is(err, ErrUserExists))
	assert.True(t, xerrors.Is(err, ErrUserExists))
}

func TestUserNotExists(t *testing.T) {
	initTestDatabase()
	defer removeTestDatabase()

	err := Auth.CreateUser(user1)
	assert.NoError(t, err)
	_, err = Auth.FindUserByLogin(user2.Login)
	assert.Error(t, err)
}

func TestUpdateUser(t *testing.T) {
	initTestDatabase()
	defer removeTestDatabase()

	err := Auth.CreateUser(user1)
	assert.NoError(t, err)
	updated := user1
	updated.FirstName = "Updated"
	err = Auth.UpdateUser(updated)
	assert.NoError(t, err)
	updated, err = Auth.FindUserByLogin(user1.Login)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.FirstName)
	assert.Equal(t, user1.LastName, updated.LastName)
	assert.Equal(t, user1.Login, updated.Login)
	assert.Equal(t, user1.PasswordHash, updated.PasswordHash)
}

func TestFindUsers(t *testing.T) {
	initTestDatabase()
	defer removeTestDatabase()

	err := Auth.CreateUser(user1)
	assert.NoError(t, err)
	err = Auth.CreateUser(user2)
	assert.NoError(t, err)
	err = Auth.CreateUser(user3)
	assert.NoError(t, err)

	users, err := Auth.FindUsers("No users", "")
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Empty(t, users)

	users, err = Auth.FindUsers(user1.Login[:1], "")
	assert.NoError(t, err)
	if !assert.Equal(t, 1, len(users)) {
		return
	}
	found := users[0]
	checkEqual(t, user1, found)

	users, err = Auth.FindUsers(user1.Login[:2], user1.Login)
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Empty(t, users)

	users, err = Auth.FindUsers(user1.FirstName[:3], "")
	assert.NoError(t, err)
	if !assert.Equal(t, 1, len(users)) {
		return
	}
	found = users[0]
	checkEqual(t, user1, found)

	users, err = Auth.FindUsers(user1.LastName[:4], "")
	assert.NoError(t, err)
	if !assert.Equal(t, 1, len(users)) {
		return
	}
	found = users[0]
	checkEqual(t, user1, found)
}
