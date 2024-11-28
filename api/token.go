package api

import (
	"errors"
	"fmt"
	"net/http"
	db "ocr/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

type ReNewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ReNewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// ReNewAccessToken
func (server *Server) ReNewAccessToken(ctx *gin.Context) {
	var req ReNewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err := errors.New("input is not valid, Please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//token verification due to invalied token and other errors
	refreshPayload, err := server.tokenMaker.VerfiyToken(req.RefreshToken)
	if err != nil {
		err := errors.New("verification not successfull, Please try again")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	//session retrival accociate with token ID, if not found returns an error
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = errors.New("session retrival failed not found, Please try again")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		err = errors.New("session retrival failed , Please try again")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	//session check, due to is session blocked returns with Unauthrized status with error
	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	//verifies the session username matches the refresh payload username
	if session.UserID != refreshPayload.UserID {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	//verifies the session refresh token matches the provided refresh token
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("session token missmatched")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	//check session has expired of not
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("exipired session")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	//After create a new token using createToken
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.UserName,
		refreshPayload.UserID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		err = errors.New("token creation failed, Please try again")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	//response containing renewed acces token and expiration time
	rsp := ReNewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
