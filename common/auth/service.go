package auth

import (
	"encoding/json"
	"strconv"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

// CreateTokenData ...
type CreateTokenData struct {
	Authentication Authentication
	ExpirationTime int // in seconds
}

// Service ...
type Service struct {
	secret string
	cache  redis.Conn
}

// NewAuthService ...
func NewAuthService(secret string, redisConn redis.Conn) *Service {
	return &Service{
		secret: secret,
		cache:  redisConn,
	}
}

// CreateToken ...
func (service *Service) CreateToken(data CreateTokenData) (string, error) {
	token := uuid.NewV4().String()

	expiry := strconv.Itoa(data.ExpirationTime)
	value, _ := json.Marshal(&data.Authentication)
	_, err := service.cache.Do("SETEX", token, expiry, string(value))

	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyToken ...
func (service *Service) VerifyToken(token string) (Authentication, bool, error) {
	var auth Authentication

	response, err := service.cache.Do("GET", token)
	if response == nil {
		// verify token is stored in cache
		// empty string means token is not found
		return Authentication{}, false, nil
	}

	value, _ := redis.Bytes(response, err)
	err = json.Unmarshal(value, &auth)

	if err != nil {
		// handle serialisation and unmarshal error
		return Authentication{}, false, err
	}

	return auth, true, nil
}
