package domain

import "github.com/iqdf/benjerry-service/common/auth"

// AuthService ...
type AuthService interface {
	CreateToken(data auth.CreateTokenData) (token string, err error)
	VerifyToken(token string) (auths auth.Authentication, success bool, err error)
}
