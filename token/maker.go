package token

import "time"

// Maker is is interface for managing tokens
type Maker interface {
	//createToken create a token for specific username and duration
	CreateToken(username string, userid int64, duration time.Duration) (string, *Payload, error)

	//verify token checks if token valid or not
	VerfiyToken(token string) (*Payload, error)
}
