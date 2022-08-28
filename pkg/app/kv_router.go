package app

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
	"github.com/kosha/redis-connector/pkg/models"
	"net/http"
	"strconv"
	"time"
)

// getAllKeys godoc
// @Summary Get all keys in redis with an optional prefix
// @Description List all keys
// @Tags KV
// @Param prefix query string false "Key Prefix"
// @Param count query integer false "Count"
// @Accept  json
// @Produce  json
// @Success 200 {array}   string
// @Router /api/v1/keys [get]
func (a *App) getAllKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	
	prefix := r.FormValue("prefix")
	if prefix == "" {
		prefix = "*"
	} else {
		prefix = prefix + ":*"
	}

	count, _ := strconv.Atoi(r.FormValue("count"))

	if count <= 0 {
		count = 0
	}

	var keys []string
	ctx := context.Background()
	iter := a.RedisClient.Scan(ctx, 0, prefix, int64(count)).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		a.Log.Errorf("Error fetching all keys with prefix: %s, count: %d", prefix, count)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, keys)

}

// getSetElements godoc
// @Summary Get all set elements in redis with an optional prefix
// @Description Get Set elements
// @Tags Sets
// @Param key path string false "Enter key"
// @Param prefix query string false "Key Prefix"
// @Param count query integer false "Count"
// @Param sort query boolean false "Sort Keys"
// @Accept  json
// @Produce  json
// @Success 200 {array}   string
// @Router /api/v1/set/{key} [get]
func (a *App) getSetElements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	prefix := r.FormValue("prefix")
	if prefix == "" {
		prefix = "*"
	} else {
		prefix = prefix + ":*"
	}

	count, _ := strconv.Atoi(r.FormValue("count"))

	if count <= 0 {
		count = 0
	}

	sort, err := strconv.ParseBool(r.FormValue("sort"))
	if err != nil {
		a.Log.Errorf("Sort query param is not a boolean. %s", sort)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	vars := mux.Vars(r)
	key := vars["key"]

	var setElements []string
	ctx := context.Background()
	var iter *redis.ScanIterator
	if !sort {
		iter = a.RedisClient.SScan(ctx, key, 0, prefix, int64(count)).Iterator()
	} else {
		iter = a.RedisClient.ZScan(ctx, key, 0, prefix, int64(count)).Iterator()

	}
	for iter.Next(ctx) {
		setElements = append(setElements, iter.Val())
	}
	if err := iter.Err(); err != nil {
		a.Log.Errorf("Error fetching all set elements key: %s with prefix: %s, count: %d", key, prefix, count)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, setElements)

}

// addElementsToASet godoc
// @Summary Add members for a set key in redis
// @Description Add members for a set
// @Tags Sets
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Param text body models.Set false "Add memebers to a set"
// @Success 200 {object}  string
// @Router /api/v1/set/{key} [post]
func (a *App) addElementsToASet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	var kv models.Set
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kv); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := context.Background()

	val, err := a.RedisClient.SAdd(ctx, key, kv.Value).Result()

	if err != nil {
		a.Log.Errorf("Error adding elements for set key %s", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}

// removeElementsFromASet godoc
// @Summary Remove members from a set key in redis
// @Description Remove members from a set
// @Tags Sets
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Param text body models.Set false "Remove memebers from a set"
// @Success 200 {object}  string
// @Router /api/v1/set/{key} [delete]
func (a *App) removeElementsFromASet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	var kv models.Set
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kv); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := context.Background()

	val, err := a.RedisClient.SRem(ctx, key, kv.Value).Result()

	if err != nil {
		a.Log.Errorf("Error removing elements from the hash key %s", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}

// getHashElements godoc
// @Summary Get all hash elements in redis with an optional prefix
// @Description Get Hash Elements
// @Tags Hash
// @Param key path string false "Enter key"
// @Param prefix query string false "Key Prefix"
// @Param count query integer false "Count"
// @Accept  json
// @Produce  json
// @Success 200 {array}   string
// @Router /api/v1/hash/{key} [get]
func (a *App) getHashElements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	prefix := r.FormValue("prefix")
	if prefix == "" {
		prefix = "*"
	} else {
		prefix = prefix + ":*"
	}

	count, _ := strconv.Atoi(r.FormValue("count"))

	if count <= 0 {
		count = 0
	}

	vars := mux.Vars(r)
	key := vars["key"]

	var hashElements []string
	ctx := context.Background()
	iter := a.RedisClient.HScan(ctx, key, 0, prefix, int64(count)).Iterator()
	for iter.Next(ctx) {
		hashElements = append(hashElements, iter.Val())
	}
	if err := iter.Err(); err != nil {
		a.Log.Errorf("Error fetching all hash elements key: %s with prefix: %s, count: %d", key, prefix, count)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, hashElements)

}

// addElementsToAHash godoc
// @Summary Add members to a hash key in redis
// @Description Add members to a hash
// @Tags Hash
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Param text body models.Hash false "Add memebers to a hash"
// @Success 200 {object}  string
// @Router /api/v1/hash/{key} [post]
func (a *App) addElementsToAHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	var kv models.Hash
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kv); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := context.Background()

	val, err := a.RedisClient.HMSet(ctx, key, kv.Value).Result()

	if err != nil {
		a.Log.Errorf("Error adding elements for hash key %s", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}

// removeElementsFromAHash godoc
// @Summary Remove members from a hash key in redis
// @Description Remove members from a hash
// @Tags Hash
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Param text body models.Hash false "Remove memebers from a hash"
// @Success 200 {object}  string
// @Router /api/v1/hash/{key} [delete]
func (a *App) removeElementsFromAHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	var kv models.Hash
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kv); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := context.Background()

	val, err := a.RedisClient.HDel(ctx, key, kv.Value...).Result()

	if err != nil {
		a.Log.Errorf("Error removing elements from the hash key %s", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}

// getSingleKey godoc
// @Summary Get single key in redis
// @Description Get key
// @Tags KV
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Success 200 {object}  string
// @Router /api/v1/kv/{key} [get]
func (a *App) getSingleKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	ctx := context.Background()
	val, err := a.RedisClient.Get(ctx, key).Result()

	if err != nil {
		a.Log.Errorf("Error fetching value for %s key", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}

// setValueForSingleKey godoc
// @Summary Set value for a key in redis
// @Description Set key value
// @Tags KV
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Param text body models.KV false "Enter value for a key"
// @Success 200 {object}  string
// @Router /api/v1/kv/{key} [post]
func (a *App) setValueForSingleKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	var kv models.KV
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kv); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := context.Background()
	var expiry time.Duration
	if kv.Expiry != 0 {
		expiry = time.Duration(kv.Expiry) * time.Second
	}
	val, err := a.RedisClient.Set(ctx, key, kv.Value, expiry).Result()

	if err != nil {
		a.Log.Errorf("Error setting value for %s key", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}

// removeSingleKey godoc
// @Summary Remove a key in redis
// @Description Remove key value
// @Tags KV
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Success 200 {object}  string
// @Router /api/v1/kv/{key} [delete]
func (a *App) removeSingleKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	ctx := context.Background()
	val, err := a.RedisClient.Del(ctx, key).Result()

	if err != nil {
		a.Log.Errorf("Error removing key %s", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}

// ping godoc
// @Summary Ping redis server
// @Description Check connectivity of a redis server
// @Tags Connectivity
// @Accept  json
// @Produce  json
// @Success 200 {object}  boolean
// @Router /api/v1/ping [get]
func (a *App) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	ctx := context.Background()
	res, err := a.RedisClient.Ping(ctx).Result()

	if err != nil {
		a.Log.Errorf("Error pinging the redis server. Error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, res)

}

// getListElements godoc
// @Summary Get all list elements in redis with an optional prefix
// @Description Get List elements
// @Tags List
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Success 200 {array}   string
// @Router /api/v1/list/{key} [get]
func (a *App) getListElements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	ctx := context.Background()
	res, err := a.RedisClient.LRange(ctx, key, 0, -1).Result()

	if err != nil {
		a.Log.Errorf("Error fetching all list elements key: %s", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, res)

}

// addElementsToAList godoc
// @Summary Add members for a list key in redis
// @Description Add members for a list
// @Tags List
// @Param key path string false "Enter key"
// @Accept  json
// @Produce  json
// @Param text body models.List false "Add memebers to a list"
// @Success 200 {object}  string
// @Router /api/v1/list/{key} [post]
func (a *App) addElementsToAList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	key := vars["key"]

	var kv models.List
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kv); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := context.Background()

	val, err := a.RedisClient.LPush(ctx, key, kv.Value).Result()

	if err != nil {
		a.Log.Errorf("Error adding elements for set key %s", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}

// removeElementsFromAList godoc
// @Summary Remove members from a list key in redis
// @Description Remove members from a list
// @Tags List
// @Param key path string false "Enter key"
// @Param count query integer false "Count"
// @Accept  json
// @Produce  json
// @Param text body models.List false "Remove memebers from a list"
// @Success 200 {object}  string
// @Router /api/v1/list/{key} [delete]
func (a *App) removeElementsFromAList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	countStr := r.FormValue("count")

	vars := mux.Vars(r)
	key := vars["key"]

	var kv models.Hash
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kv); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := context.Background()

	var val int64
	var err error
	if countStr != "" {
		c, err := strconv.Atoi(countStr)
		if err != nil {
			a.Log.Errorf("Count not a valid integer", err)
			respondWithError(w, http.StatusBadRequest, "Invalid count parameter. Should be either <0, 0, >0")
			return
		}
		val, err = a.RedisClient.LRem(ctx, key, int64(c), kv.Value).Result()
	} else {
		val, err = a.RedisClient.LRem(ctx, key, -1, kv.Value).Result()
	}

	if err != nil {
		a.Log.Errorf("Error removing elements from the hash key %s", key)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, val)

}
