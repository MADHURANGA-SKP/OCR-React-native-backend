package api

import (
	"fmt"
	db "ocr/db/sqlc"
	"ocr/token"
	"ocr/util"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// server serves http requests
// holds containing configuration details, store database connection,
// responsible for generating and verifying access tokens
// holds core routing component for handling HTTP requests
// holds instance of worker.TaskDistributor to distribute tasks to worker processes
type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

// NewServer create a http server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	//calls the setupRouter method to configure routing for the server
	server.setupRouter()
	//returns created Server and is there any potential error
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	//Cross-Origin Resource Sharing (CORS) uses the cors.New function.
	//And it allows requests from specific origins to interact with the server's API.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://*", "https://*", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := router.Group("/ocr")
	{
		v1.POST("/signup", server.CreateUser)
		v1.POST("/login", server.LoginUser)
		v1.POST("/refresh/accesstoken", server.ReNewAccessToken)
		v1.GET("/getuser", server.GetUser)
		v1.GET("/getusers", server.GetImageConversion)
		v1.DELETE("/deluser", server.DeleteUser)
		v1.PATCH("/updateuser", server.UpadteUser)
		v1.POST("/imageconv", server.CreateImageConversion)
	}

	server.router = router
}

// start runs the http server on a specified address.
// Start defined in the main.go
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// This creates a simple JSON response containing an "error" key with the returning error message as value
// And it can be called as errorResponse(err)
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// This creates a simple JSON response containing an "msg" key with the returning message as value
// And it can be called as errorResponse(err)
func messageResponse(msg string) gin.H {
	return gin.H{"message": msg}
}
