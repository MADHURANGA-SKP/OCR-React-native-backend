package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Defines the minimum  size for the secret key
const minSecretKeySize = 32

// create struct of JWT json web token maker
type JWTMaker struct {
	secretkey string
}

// NewJWTMaker createts a new JWTMaker
func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, fmt.Errorf("invalied key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretkey}, nil
}

// CreateToken creates  a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, userid int64, duration time.Duration) (string, *Payload, error) {
	//create new payload for JWT token using username,role,userid,duration, if is there no error return payload, and error occured, return the error
	payload, err := Newpayload(username, userid, duration)
	if err != nil {
		return "", payload, err
	}

	//create new JWT token using payload and signing method
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	//sign the token with secret key and return the created token, if is there no errors
	token, err := jwtToken.SignedString([]byte(maker.secretkey))
	return token, payload, err
}

// VerifyToken checks if token is valied or not
func (maker *JWTMaker) VerfiyToken(token string) (*Payload, error) {
	//define key funtion that validate the token
	keyFunc := func(token *jwt.Token) (interface{}, error) {

		//check the sign in method valid or not
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		//return the key for validation
		return []byte(maker.secretkey), nil
	}

	//parse the token with payload and keyfuntion its meets the included details and JWT funtionalitis
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {

		//validate errors
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	//get validated token
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	//return the validated token
	return payload, nil
}
