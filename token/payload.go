package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Diffrent types of errors are returns by verify token funtion
var (
	ErrInvalidToken = errors.New("token is invalied")
	ErrExpiredToken = errors.New("token is expired")
)

// Payload contains the payload data of the token that generate for each user
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"user_name"`
	UserID    int64     `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Newpayload creates a new token with a specific username and duration
func Newpayload(username string, userid int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserName:  username,
		UserID:    userid,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if token payload valied or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
