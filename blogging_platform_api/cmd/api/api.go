package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-projects/blogging_platform_api/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := gin.Default()

	blogStore := internal.NewStore(s.db)
	blogHandler := internal.NewHandler(blogStore)
	subrouter := router.Group("/api/v1")
	{
		blogHandler.BlogRoutes(subrouter)
	}

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
