package api

import (
	"log/slog"
	"os"

	"github.com/cv711/odin-takehome/server/db"
	"github.com/gin-gonic/gin"
)

type API struct {
	log           *slog.Logger
	db            *db.Queries
	adminPassword string
}

func NewAPI(log *slog.Logger, db *db.Queries) *API {
	return &API{
		log:           log,
		db:            db,
		adminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}

func (a *API) setupRoutes(router *gin.Engine) *gin.Engine {
	apiRouter := router.Group("/api")
	apiRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	apiRouter.POST("/auth", a.auth)
	apiRouter.POST("/signup", a.signup)
	apiRouter.GET("/user", a.authRoute, a.getUser)

	return router
}

func (a *API) Serve() {
	a.log.Info("Server starting...")
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(gin.Recovery())

	server = a.setupRoutes(server)

	port, portSet := os.LookupEnv("PORT")
	if !portSet {
		port = "8080"
	}

	a.log.Info("Listening on port " + port)
	if err := server.Run(":" + port); err != nil {
		panic(err)
	}
}
