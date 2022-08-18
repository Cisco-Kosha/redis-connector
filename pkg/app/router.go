package app

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/kosha/redis-connector/pkg/models"
	"net/http"
)

// listConnectorSpecification godoc
// @Summary Get details of the connector specification
// @Description Get all environment variables that need to be supplied
// @Tags specification
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Router /api/v1/specification/list [get]
func (a *App) listConnectorSpecification(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	respondWithJSON(w, http.StatusOK, map[string]string{
		"USERNAME":   "Redis Username",
		"PASSWORD":   "Redis Password",
		"DATABASE":   "Redis Database",
		"REDIS_HOST": "Redis Host",
		"REDIS_PORT": "Redis Port",
	})
}

// testConnectorSpecification godoc
// @Summary Test if username, password and redis host url against the specification
// @Description Check if redis server connection is successful
// @Tags specification
// @Accept  json
// @Produce  json
// @Param text body models.Specification false "Enter username, password and redis host"
// @Success 200 {object} boolean
// @Router /api/v1/specification/test [post]
func (a *App) testConnectorSpecification(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	if (*r).Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	var s models.Specification
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     s.RedisHost + ":" + s.RedisPort,
		Password: s.Password,
		DB:       0, // use default DB
		Username: s.Username,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		respondWithJSON(w, http.StatusOK, false)
	} else {
		respondWithJSON(w, http.StatusOK, true)
	}
}
