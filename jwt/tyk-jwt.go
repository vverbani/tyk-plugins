package main

import (
	//"gopkg.in/square/go-jose.v2/jwt"

	"net/http"
	"os"

	"github.com/TykTechnologies/tyk/log"
	"gopkg.in/square/go-jose.v2/json"
	//"encoding/json"
)

type jwtConfig struct {
	jwtCertFile   string   `json:"jwtCertFile"`
	jwtKeyFile    string   `json:"jwtKeyFile"`
	jwksCertFiles []string `json:"jwksCertFiles"`
}

var logger = log.Get()

func init() {
	configFileName := "/opt/tyk-plugins/tyk-jwt.json"
	var config jwtConfig
	logger.Info("jtw-plugin: Loading config file: ", configFileName)
	configFile, err := os.Open(configFileName)
	if err != nil {
		logger.Error("jtw-plugin: Cannot load ", configFileName, " : ", err)
		os.Exit(1)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		logger.Warning("Parsing ", configFile, " gave error: ", err)
	} else {
		logger.Info("jwtCertFile is: ", config.jwtCertFile)
		logger.Info("jwtKeyFile is: ", config.jwtKeyFile)
		logger.Info("jwksCertFiles are: ", config.jwksCertFiles)
	}
}

// AddFooBarHeader adds custom "Foo: Bar" header to the request
func AddFooBarHeader(rw http.ResponseWriter, r *http.Request) {
	logger.Info("Processing HTTP request in Golang plugin!!")
	r.Header.Add("Foo", "Bar")
}

func main() {}
