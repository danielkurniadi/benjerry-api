package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/iqdf/benjerry-service/common/auth"
	"github.com/iqdf/benjerry-service/domain"
)

// UserHandler ...
type UserHandler struct {
	service       domain.UserService
	authService   *auth.Service
	sessionExpiry time.Duration
}

// NewUserHandler ...
func NewUserHandler(service domain.UserService, authService *auth.Service, sessionExpiry time.Duration) *UserHandler {
	return &UserHandler{
		service:       service,
		authService:   authService,
		sessionExpiry: sessionExpiry,
	}
}

// Routes register handle func with the path url
func (handler *UserHandler) Routes(router *mux.Router) {
	// Register handler methods to router here...
	router.Handle("/login", handler.handleLogin()).Methods("POST")
	router.Handle("/signup", handler.handleSignUp()).Methods("POST")
}

func (handler *UserHandler) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		ctx := r.Context()
		username, rawpass, ok := r.BasicAuth()

		if !ok {
			failAuthentication(w)
			return
		}

		user, err := handler.service.LoginUser(ctx, username, rawpass)

		if err != nil {
			failAuthentication(w)
			return
		}

		expiry := int(handler.sessionExpiry.Seconds())
		authentication := auth.Authentication{
			ID:             username,
			Authorizations: user.Authorizations,
		}
		createTokenData := auth.CreateTokenData{
			Authentication: authentication,
			ExpirationTime: expiry,
		}

		sessionToken, err := handler.authService.CreateToken(createTokenData)

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token", // TODO: move to const
			Value:   sessionToken,
			Expires: time.Now().Add(480 * time.Second),
		})

		w.WriteHeader(200)
		w.Write([]byte("Login success"))
	}
}

func (handler *UserHandler) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		ctx := r.Context()
		username, rawpass, ok := r.BasicAuth()

		if !ok {
			failAuthentication(w)
			return
		}

		err := handler.service.RegisterUser(ctx, username, rawpass, false)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("signup: Internal Server error.\n"))
			return
		}

		w.WriteHeader(201)
		w.Write([]byte("Account created successfully"))
	}
}

func failAuthentication(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="User Visible Realm`)
	w.WriteHeader(401)
	w.Write([]byte("Unauthorized.\n"))
}
