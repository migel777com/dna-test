package server

import (
	h "dna-test/api/handler"
	"dna-test/api/middleware"
	"dna-test/config"
	"dna-test/models"
	srv "dna-test/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Server struct {
	Configuration *config.Config
	Framework     *gin.Engine
	Db            models.DbClient
	Cache         models.CacheClient
}

func NewApiServer(config *config.Config, db models.DbClient, cache models.CacheClient) *Server {
	return &Server{
		Configuration: config,
		Framework:     gin.Default(),
		Db:            db,
		Cache:         cache,
	}
}

func (s *Server) Init() {
	s.Framework.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}),
		middleware.JSONMiddleware(),
		middleware.ErrorHandler())

	service := srv.NewService(s.Db, s.Cache)

	h.BindRoutes(service, s.Framework)
}
