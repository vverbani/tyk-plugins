package main

import (
	"gopkg.in/square/go-jose.v2/jwt"
	"github.com/TykTechnologies/tyk/log"
	"github.com/tkanos/gonfig"
	"encoding/json"
	"os"
)

type jwtConfig struct {
	jwtCertFile  string
	jwtKeyFile   string
	jwksCertFiles []string
}

var logger = log.Get()

func init () {
	configFile := "/opt/tyk-plugin/tyk-jwt.json"
	logger.Info("jtw-plugin: Loading config file: %s", configFile)
	/* f, err := os.Open(configFile)
	if err != nil {
		logger.Error("jtw-plugin: Cannot load %s: %s"m configFile, err)
		exit
	} */
	config := jwtConfig
	err := gonfig.GetConf(configFile, &config)
}

func main() {}