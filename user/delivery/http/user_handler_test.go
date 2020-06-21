package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/iqdf/benjerry-service/domain"
	"github.com/iqdf/benjerry-service/domain/mocks"
	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

var (
	contextType   = mock.Anything
	usernameType  = mock.AnythingOfType("string")
	rawpassType   = mock.AnythingOfType("string")
	tokenDataType = mock.AnythingOfType("auth.CreateTokenData")
	isAdminType   = mock.AnythingOfType("bool")
)

func TestHandleLoginSuccess(t *testing.T) {
	userService := new(mocks.UserService)
	authService := new(mocks.AuthService)

	username, hashpass := "usertest", createMockHashPassword()
	rawpass := "passwordtest"
	token := createMockToken()
	mockUser := createMockUser(username, hashpass)

	userService.
		On("LoginUser", contextType, usernameType, rawpassType).
		Return(mockUser, nil).
		Once()

	authService.
		On("CreateToken", tokenDataType).
		Return(token, nil).
		Once()

	request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(""))
	request.SetBasicAuth(username, rawpass)

	recorder := httptest.NewRecorder()

	userHandler := NewUserHandler(userService, authService, 640*time.Second)
	loginHandler := userHandler.handleLogin()

	loginHandler(recorder, request)
	assert.Equal(t, recorder.Code, 200)
	assert.Equal(t, recorder.Body.String(), "Login success\n")
}

func TestHandleLoginNoAuth(t *testing.T) {
	userService := new(mocks.UserService)
	authService := new(mocks.AuthService)

	token := createMockToken()
	mockUser := createMockUser("usertest12", createMockHashPassword())

	userService.
		On("LoginUser", contextType, usernameType, rawpassType).
		Return(mockUser, nil).
		Once()

	authService.
		On("CreateToken", tokenDataType).
		Return(token, nil).
		Once()

	request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(""))
	recorder := httptest.NewRecorder()

	userHandler := NewUserHandler(userService, authService, 640*time.Second)
	loginHandler := userHandler.handleLogin()

	loginHandler(recorder, request)
	assert.Equal(t, recorder.Code, 401)
	assert.Equal(t, recorder.Body.String(), "Unauthorized.\n")
}

func TestHandleLoginBadUsername(t *testing.T) {
	userService := new(mocks.UserService)
	authService := new(mocks.AuthService)

	badusername := "*912(!#$@a"
	password := "passwordtest"
	token := createMockToken()
	mockUser := createMockUser(badusername, createMockHashPassword())

	userService.
		On("LoginUser", contextType, usernameType, rawpassType).
		Return(mockUser, nil).
		Once()

	authService.
		On("CreateToken", tokenDataType).
		Return(token, nil).
		Once()

	request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(""))
	request.SetBasicAuth(badusername, password)

	recorder := httptest.NewRecorder()

	userHandler := NewUserHandler(userService, authService, 640*time.Second)
	loginHandler := userHandler.handleLogin()

	loginHandler(recorder, request)
	assert.Equal(t, recorder.Code, 400)
}

func TestHandleLoginAuthFail(t *testing.T) {
	userService := new(mocks.UserService)
	authService := new(mocks.AuthService)

	badusername := "usertest123"
	password := "passwordtest"
	token := createMockToken()

	userService.
		On("LoginUser", contextType, usernameType, rawpassType).
		Return(domain.User{}, domain.ErrAuthFail).
		Once()

	authService.
		On("CreateToken", tokenDataType).
		Return(token, nil).
		Once()

	request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(""))
	request.SetBasicAuth(badusername, password)

	recorder := httptest.NewRecorder()

	userHandler := NewUserHandler(userService, authService, 640*time.Second)
	loginHandle := userHandler.handleLogin()

	loginHandle(recorder, request)
	assert.Equal(t, recorder.Code, 401)
}

func TestHandleSignUpSuccess(t *testing.T) {
	userService := new(mocks.UserService)
	authService := new(mocks.AuthService)

	username, password := "usertest", "passwordtest"

	userService.
		On("RegisterUser", contextType, usernameType, rawpassType, isAdminType).
		Return(nil)

	request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(""))
	request.SetBasicAuth(username, password)

	recorder := httptest.NewRecorder()

	userHandler := NewUserHandler(userService, authService, 640*time.Second)
	signupHandle := userHandler.handleSignUp()

	signupHandle(recorder, request)
	assert.Equal(t, recorder.Code, 201)
	assert.Equal(t, recorder.Body.String(), "Account created successfully\n")
}

func TestHandleSignUpBadUsername(t *testing.T) {
	userService := new(mocks.UserService)
	authService := new(mocks.AuthService)

	badusername, password := "*912(!#$@a", "passwordtest"

	request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(""))
	request.SetBasicAuth(badusername, password)

	recorder := httptest.NewRecorder()

	userHandler := NewUserHandler(userService, authService, 640*time.Second)
	signupHandle := userHandler.handleSignUp()

	signupHandle(recorder, request)
	assert.Equal(t, recorder.Code, 400)
}

func TestHandleSignUpBadPassword(t *testing.T) {
	userService := new(mocks.UserService)
	authService := new(mocks.AuthService)

	username, badpassword := "usertest", "12a3"

	request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(""))
	request.SetBasicAuth(username, badpassword)

	recorder := httptest.NewRecorder()

	userHandler := NewUserHandler(userService, authService, 640*time.Second)
	signupHandle := userHandler.handleSignUp()

	signupHandle(recorder, request)
	assert.Equal(t, recorder.Code, 400)
}

func TestHandleSignUpConflict(t *testing.T) {
	userService := new(mocks.UserService)
	authService := new(mocks.AuthService)

	username, badpassword := "usertest", "passwordtest"

	userService.
		On("RegisterUser", contextType, usernameType, rawpassType, isAdminType).
		Return(domain.ErrConflict)

	request, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(""))
	request.SetBasicAuth(username, badpassword)

	recorder := httptest.NewRecorder()

	userHandler := NewUserHandler(userService, authService, 640*time.Second)
	signupHandle := userHandler.handleSignUp()

	signupHandle(recorder, request)
	assert.Equal(t, recorder.Code, 200)
	assert.Equal(t, recorder.Body.String(), "User with same username already exist.\n")
}

func createMockUser(username, hashpassword string) domain.User {
	return domain.User{
		Username:     username,
		HashPassword: hashpassword,
	}
}

func createMockToken() string {
	return "a70e8544-a708-4372-8ef9-39b1d3d6e3fe"
}

func createMockHashPassword() string {
	return "$2a$04$5cHVB3Jgu4b6GDpCuW8zGu/jmAuDVepej.aW7oQWlksFsOOFuqlTO"
}
