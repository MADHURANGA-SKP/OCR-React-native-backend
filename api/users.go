package api

import (
	"errors"
	"net/http"
	db "ocr/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgtype"
)

// create useer
type CreateUserRequest struct {
	UserName       string `json:"user_name" binding:"required"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email" binding:"required,email"`
	HashedPassword string `json:"hashed_password" binding:"required,min=8"`
}

type CreateUserResponse struct {
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) CreateUserResponse {
	return CreateUserResponse{
		UserID:    user.UserID,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}

func (server Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(req); err != nil {
		err = errors.New("input is not valid, Please Check")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUsersParams{
		UserName:       req.UserName,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.LastName,
		HashedPassword: req.HashedPassword,
	}

	user, err := server.store.CreateUsers(ctx, arg)
	if err != nil {
		err = errors.New("failed to create course, Please try again")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

// get user
type GetUserRequest struct {
	UserID int64 `json:"user_id"`
}

type GetUserResponse struct {
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (server Server) GetUserRequest(ctx *gin.Context) {
	var req GetUserRequest

	if err := ctx.ShouldBindJSON(req); err != nil {
		err = errors.New("user not found, Please enter your correct User Name")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserID)
	if err != nil {
		err = errors.New("user not found, Please enter your correct User Name")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Delete user
type DeleteUserRequest struct {
	UserID int64 `json:"user_id"`
}

func (server Server) DeleteUserRequest(ctx *gin.Context) {
	var req DeleteUserRequest

	if err := ctx.ShouldBindJSON(req); err != nil {
		err = errors.New("user not found, Please enter your correct User Name")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err := server.store.DeleteUsers(ctx, req.UserID)
	if err != nil {
		err = errors.New("user not found, Please enter your correct User Name")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	response := messageResponse("Failed to delete User")
	ctx.JSON(http.StatusOK, response)
}

// Upadte user
type UpadteUserRequest struct {
	UserID         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type UpadteUserResponse struct {
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (server Server) UpadteUserRequest(ctx *gin.Context) {
	var req UpadteUserRequest

	if err := ctx.ShouldBindJSON(req); err != nil {
		err = errors.New("user not found, Please enter your correct User Name")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.UpdateUsersParams{
		UserID: req.UserID,
		UserName: pgtype.Text{
			String: req.UserName,
			Valid:  true,
		},
		FirstName: pgtype.Text{
			String: req.FirstName,
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: req.LastName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: req.Email,
			Valid:  true,
		},
		HashedPassword: pgtype.Text{
			String: req.HashedPassword,
			Valid:  true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	user, err := server.store.UpdateUsers(ctx, arg)
	if err != nil {
		err = errors.New("user not found, Please enter your correct User Name")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
