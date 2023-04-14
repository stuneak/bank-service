package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stuneak/simplebank/sqlc_internal"
)

type Server struct {
	store  sqlc_internal.Store
	router *gin.Engine
}

func NewServer(store sqlc_internal.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"Error": err.Error()}
}
