package main

import (
	"net/http"

	logger "github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/storage"
)

var log = logger.Get()

const pluginDefaultKeyPrefix = "PetesPlugin-data:"

func tykStoreData(key, value string) {
	ttl := int64(1000)
	store := storage.RedisCluster{KeyPrefix: pluginDefaultKeyPrefix}
	store.SetKey(key, value, ttl)
}

func tykGetData(key string) string {
	store := storage.RedisCluster{KeyPrefix: pluginDefaultKeyPrefix}
	val, _ := store.GetKey(key)
	return val
}

// CallRedis based on Policy ID
func CallRedis(rw http.ResponseWriter, r *http.Request) {

	log.Info("Start CallRedis Plugin")

	tykStoreData("Pete", "Woz Here")
	log.Info("CallRedis: saved")
	log.Info("CallRedis: loaded: " + tykGetData("Pete"))

	log.Info("End CallRedis Plugin")

}

func main() {}
