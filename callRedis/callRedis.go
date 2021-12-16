package main

import (
	"net/http"

	logger "github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/storage"
)

const pluginDefaultKeyPrefix = "PetesPlugin-data:"

var log = logger.Get()
var store = storage.RedisCluster{KeyPrefix: pluginDefaultKeyPrefix}

func tykStoreData(key, value string) {
	ttl := int64(1000)
	store.SetKey(key, value, ttl)
}

func tykGetData(key string) string {
	val, _ := store.GetKey(key)
	return val
}

// CallRedis to poke stuff in there for funzies
func CallRedis(rw http.ResponseWriter, r *http.Request) {

	log.Info("Start CallRedis Plugin")

	tykStoreData("Pete", "Woz Here")
	log.Info("CallRedis: saved to redis")
	log.Info("CallRedis: value retrieved: " + tykGetData("Pete"))

	log.Info("End CallRedis Plugin")

}

func main() {}
