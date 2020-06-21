package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/iqdf/benjerry-service/common/auth"
	validatorLib "github.com/iqdf/benjerry-service/common/validator"
	"github.com/iqdf/benjerry-service/domain"
)

// UserHandler ...
type UserHandler struct {
	userService   domain.UserService
	authService   domain.AuthService
	sessionExpiry time.Duration
}

// NewUserHandler ...
func NewUserHandler(
	service domain.UserService,
	authService domain.AuthService,
	sessionExpiry time.Duration,
) *UserHandler {
	return &UserHandler{
		userService:   service,
		authService:   authService,
		sessionExpiry: sessionExpiry,
	}
}

// Routes register handle func with the path url
func (handler *UserHandler) Routes(router *mux.Router) {
	// Register handler methods to router here...
	router.Handle("/login", handler.handleLogin()).Methods("POST")
	router.Handle("/signup", handler.handleSignUp()).Methods("POST")
	router.Handle("/admin", handler.handleSignUpAdmin()).Methods("POST")
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

		if errs := validatorLib.ValidateVar(username, "min=3,max=20,alphanum"); errs != nil {
			failBadCredentialParams(w, errs)
			return
		}

		if errs := validatorLib.ValidateVar(rawpass, "min=8,max=30,ascii"); errs != nil {
			failBadCredentialParams(w, errs)
			return
		}

		user, err := handler.userService.LoginUser(ctx, username, rawpass)

		if err == domain.ErrAuthFail || err == domain.ErrResourceNotFound {
			failAuthentication(w)
			return
		}

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("login: internal server error: " + err.Error()))
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
		w.Write([]byte("Login success\n"))
	}
}

func (handler *UserHandler) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		ctx := r.Context()
		username, rawpass, ok := r.BasicAuth()

		if !ok {
			failBadCredentialParams(w, domain.ErrBadParamInput)
			return
		}

		if errs := validatorLib.ValidateVar(username, "min=3,max=20,alphanum"); errs != nil {
			failBadCredentialParams(w, errs)
			return
		}

		if errs := validatorLib.ValidateVar(rawpass, "min=8,max=30,ascii"); errs != nil {
			failBadCredentialParams(w, errs)
			return
		}

		err := handler.userService.RegisterUser(ctx, username, rawpass, false)

		if err == domain.ErrConflict {
			w.WriteHeader(200)
			w.Write([]byte("User with same username already exist.\n"))
			return
		}

		w.WriteHeader(201)
		w.Write([]byte("Account created successfully\n"))
	}
}

func (handler *UserHandler) handleSignUpAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		ctx := r.Context()
		username, rawpass, ok := r.BasicAuth()

		if !ok {
			failBadCredentialParams(w, domain.ErrBadParamInput)
			return
		}

		if errs := validatorLib.ValidateVar(username, "min=3,max=20,alphanum"); errs != nil {
			failBadCredentialParams(w, errs)
			return
		}

		if errs := validatorLib.ValidateVar(rawpass, "min=8,max=30,ascii"); errs != nil {
			failBadCredentialParams(w, errs)
			return
		}

		err := handler.userService.RegisterUser(ctx, username, rawpass, true)

		if err == domain.ErrConflict {
			w.WriteHeader(200)
			w.Write([]byte("User with same username already exist.\n"))
			return
		}

		w.WriteHeader(201)
		w.Write([]byte("Account created successfully\n"))
	}
}

func failAuthentication(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="User Visible Realm`)
	w.WriteHeader(401)
	w.Write([]byte("Unauthorized.\n"))
}

func failBadCredentialParams(w http.ResponseWriter, errs error) {
	var message string = "login: Invalid input for username or password fields.\n"
	if verr, ok := errs.(*validatorLib.ValidationError); ok {
		message = verr.Message()
	}
	w.WriteHeader(400)
	w.Write([]byte(message))
}
