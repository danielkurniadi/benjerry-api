package service

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/iqdf/benjerry-service/domain"
	"github.com/iqdf/benjerry-service/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	appName      = "TestApp"
	contextType  = mock.Anything
	usernameType = mock.AnythingOfType("string")
	userType     = mock.AnythingOfType("domain.User")
)

func TestRegisterUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("RegisterUser-success-notAdmin", func(t *testing.T) {
		username, password, isAdmin := "usertest", "passwordtest", false
		mockUserRepo.
			On("Get", contextType, usernameType).
			Return(domain.User{}, domain.ErrResourceNotFound).
			Once()

		mockUserRepo.
			On("Create", contextType, userType).
			Return(nil).
			Once()

		userService := NewUserService(appName, mockUserRepo)
		err := userService.RegisterUser(context.TODO(), username, password, isAdmin)

		assert.NoError(t, err)
	})

	t.Run("RegisterUser-success-admin", func(t *testing.T) {
		username, password, isAdmin := "usertest", "passwordtest", true
		mockUserRepo.
			On("Get", contextType, usernameType).
			Return(domain.User{}, domain.ErrResourceNotFound).
			Once()
	
		mockUserRepo.
			On("Create", contextType, userType).
			Return(nil).
			Once()

		userService := NewUserService(appName, mockUserRepo)
		err := userService.RegisterUser(context.TODO(), username, password, isAdmin)

		assert.NoError(t, err)
	})

	t.Run("RegisterUser-failed", func(t *testing.T) {
		var dbErr = domain.ErrConflict
		username, password, isAdmin := "usertest", "passwordtest", false
		mockUserRepo.
			On("Get", contextType, usernameType).
			Return(domain.User{Username: username}, nil).
			Once()
	
			mockUserRepo.
			On("Create", contextType, userType).
			Return(dbErr).
			Once()

		userService := NewUserService(appName, mockUserRepo)
		err := userService.RegisterUser(context.TODO(), username, password, isAdmin)

		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}

func TestLoginUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("LoginUser-success", func(t *testing.T) {
		username, password := "usertest", "passwordtest"
		mockUser := createMockUser(username, password)

		mockUserRepo.
			On("Get", contextType, usernameType).
			Return(mockUser, nil).
			Once()

		userService := NewUserService(appName, mockUserRepo)
		user, err := userService.LoginUser(context.TODO(), username, password)

		assert.NoError(t, err)
		assert.True(t, cmp.Equal(user, mockUser))
	})

	t.Run("LoginUser-wrongpass", func(t *testing.T) {
		username, password := "usertest", "passwordtest"
		mockUser := createMockUser(username, password)

		mockUserRepo.
			On("Get", contextType, usernameType).
			Return(mockUser, nil).
			Once()

		userService := NewUserService(appName, mockUserRepo)
		user, err := userService.LoginUser(context.TODO(), username, "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, err, domain.ErrAuthFail)
		assert.True(t, cmp.Equal(user, domain.User{}))
	})

	t.Run("LoginUser-notfound", func(t *testing.T) {
		username, password := "usertest", "passwordtest"

		mockUserRepo.
			On("Get", contextType, usernameType).
			Return(domain.User{}, domain.ErrResourceNotFound).
			Once()

		userService := NewUserService(appName, mockUserRepo)
		user, err := userService.LoginUser(context.TODO(), username, password)

		assert.Error(t, err)
		assert.Equal(t, err, domain.ErrResourceNotFound)
		assert.True(t, cmp.Equal(user, domain.User{}))
	})
}

func createMockUser(username, password string) domain.User {
	bpassword := []byte(password)
	hashpass := hashAndSalt(bpassword)

	return domain.User{
		Username:     username,
		HashPassword: hashpass,
	}
}
