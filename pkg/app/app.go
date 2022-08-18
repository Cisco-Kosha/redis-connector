package app

import (
	"github.com/go-redis/redis/v9"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kosha/redis-connector/pkg/config"
	"github.com/kosha/redis-connector/pkg/logger"
)

type App struct {
	Router      *mux.Router
	Log         logger.Logger
	Cfg         *config.Config
	RedisClient *redis.Client
}

func router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}

// Initialize creates the necessary scaffolding of the app
func (a *App) Initialize(log logger.Logger) {

	cfg := config.Get()

	a.Cfg = cfg
	a.Log = log

	a.Router = router()

	a.RedisClient = cfg.GetRedisClient()
	a.initializeRoutes()

}

// Run starts the app and serves on the specified addr
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
