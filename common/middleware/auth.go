package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/iqdf/benjerry-service/common/auth"
	"github.com/iqdf/benjerry-service/common/consts/role"
)

// AuthKey to get authentication context from request
type AuthKey string

const authKey AuthKey = "authentication"

// AuthMiddleWare ...
func AuthMiddleWare(service *auth.Service) alice.Constructor {
	verifyAuthenticated := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve session token from request cookie
			cookie, err := r.Cookie("session_token")
			if err != nil {
				if err == http.ErrNoCookie {
					// if cookie not set, return unauthorized status
					fmt.Println("No Cookie?")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			sessionToken := cookie.Value

			// Retrieve credential from cache and verify
			auth, verified, err := service.VerifyToken(sessionToken)

			if !verified {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Session timeout. Please relogin"))
				return
			}

			fmt.Println("inserting auth to context")
			ctx := context.WithValue(r.Context(), authKey, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	return verifyAuthenticated
}

// RoleMiddleWare ...
func RoleMiddleWare(appName string) alice.Constructor {
	getRequiredRole := func(routeName string) string {
		switch {
		case strings.HasSuffix(routeName, "GET") || strings.HasSuffix(routeName, "FETCH"):
			return role.ReadPermission
		case strings.HasSuffix(routeName, "CREATE") || strings.HasSuffix(routeName, "UPDATE"):
			return role.WritePermission
		case strings.HasSuffix(routeName, "DELETE"):
			return role.DeletePermission
		}
		return role.Unauthorized
	}

	verifyFromContext := func(ctx context.Context, k AuthKey, requiredRole string) bool {
		if requiredRole == role.Unauthorized {
			return false
		}

		if v := ctx.Value(authKey); v == nil {
			// authorization context isn't set
			fmt.Println("ctx auth not set?", v)
			return false
		} else if auth, ok := v.(auth.Authentication); !ok {
			// authorization context isn't set
			fmt.Println("cannot cast auth object")
			return false
		} else if len(auth.Authorizations) > 0 {
			// check authorization roles here
			// and return true only if any role matches
			for _, a := range auth.Authorizations {
				if a.AppName == appName && a.Role == requiredRole {
					return true
				}
			}
			fmt.Println("no authorization role matching")
		}
		return false
	}

	verifyAuthorized := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Role Middleware")
			ctx := r.Context()
			var role, routeName string

			if route := mux.CurrentRoute(r); route != nil {
				routeName = route.GetName()
			}

			role = getRequiredRole(routeName)

			fmt.Println(role, routeName)
			if verifyFromContext(ctx, authKey, role) == false {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Operation not permitted"))
				return // important!
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	return verifyAuthorized
}
