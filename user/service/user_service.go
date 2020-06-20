package service

import (
	"context"
	"log"
	"time"

	"github.com/iqdf/benjerry-service/common/auth"
	"github.com/iqdf/benjerry-service/common/consts/role"
	"github.com/iqdf/benjerry-service/domain"
	"golang.org/x/crypto/bcrypt"
)

const timeout = time.Second * 10

// UserService ...
type UserService struct {
	appName  string
	userRepo domain.UserRepository
}

// NewUserService creates new service
// that provides use cases for user data/resource
func NewUserService(appName string, userRepo domain.UserRepository) *UserService {
	return &UserService{
		appName:  appName,
		userRepo: userRepo,
	}
}

// LoginUser ...
func (service *UserService) LoginUser(ctx context.Context, username, rawpass string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	user, err := service.userRepo.Get(ctx, username)
	if err == domain.ErrResourceNotFound {
		log.Printf("db notfound | username %s, hash %s, raw%s", username, user.HashPassword, rawpass)
		return domain.User{}, err
	}

	if !comparePasswords(user.HashPassword, []byte(rawpass)) {
		log.Printf("cmp password | username %s, hash %s, raw %s", user.Username, user.HashPassword, rawpass)
		return domain.User{}, domain.ErrAuthFail
	}

	return user, nil
}

// RegisterUser ...
func (service *UserService) RegisterUser(ctx context.Context, username, rawpass string, isAdmin bool) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := service.userRepo.Get(ctx, username)
	if err == nil && err != domain.ErrResourceNotFound {
		return domain.ErrConflict
	}

	var (
		user           domain.User
		authorizations []auth.Authorization
	)

	if isAdmin {
		authorizations = []auth.Authorization{
			{AppName: service.appName, Role: role.ReadPermission},
			{AppName: service.appName, Role: role.WritePermission},
			{AppName: service.appName, Role: role.DeletePermission},
		}
	} else {
		authorizations = []auth.Authorization{
			{AppName: service.appName, Role: role.ReadPermission},
		}
	}

	hashpass := hashAndSalt([]byte(rawpass))
	user = domain.User{
		Username:       username,
		HashPassword:   hashpass,
		Authorizations: authorizations,
	}
	err = service.userRepo.Create(ctx, user)
	return err
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
