package main

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/TykTechnologies/tyk/log"
	jwt "github.com/dgrijalva/jwt-go"
)

type jwtConfig struct {
	JwtCertFile   string   `json:"jwtCertFile"`
	JwtKeyFile    string   `json:"jwtKeyFile"`
	JwksCertFiles []string `json:"jwksCertFiles"`
}

var (
	jwtPublicCert *rsa.PublicKey
	jwtPrivateKey *rsa.PrivateKey
	logger        = log.Get()
)

func init() {
	configFileName := "/opt/tyk-plugins/tyk-jwt.json"

	// read the config file
	var config jwtConfig
	logger.Info("jtw-plugin: Loading config file: ", configFileName)
	configFile, err := os.Open(configFileName)
	logErr("Error", "jtw-plugin: Cannot load "+configFileName+": ", err)
	defer configFile.Close()
	// Parse the config file
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	logErr("Error", "Parsing "+configFileName+" gave error: ", err)

	// report on the values loaded
	logger.Info("jwtCertFile is: ", config.JwtCertFile)
	logger.Info("jwtKeyFile is: ", config.JwtKeyFile)
	for i := range config.JwksCertFiles {
		logger.Info("jwksCertFiles are: ", config.JwksCertFiles[i])
	}

	// start of jwt stuff
	//token := jwt.New(jwt.SigningMethodRS256)
	jwtKeyBytes, err := ioutil.ReadFile(config.JwtKeyFile)
	logErr("Error", "Opening JwtKeyFile "+config.JwtKeyFile+" failed: ", err)
	jwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(jwtKeyBytes)
	logErr("Error", "ParseRSAPrivateKeyFromPEM failed with error: ", err)
	logger.Info("jwtPrivateKey is: ", jwtPrivateKey) // lets not log our private key

	jwtCertBytes, err := ioutil.ReadFile(config.JwtCertFile)
	logErr("Error", "Opening JwtKeyFile "+config.JwtCertFile+" failed: ", err)
	jwtPublicCert, err = jwt.ParseRSAPublicKeyFromPEM(jwtCertBytes)
	logErr("Error", "ParseRSAPublicKeyFromPEM failed with error: ", err)
	logger.Info("jwtPublicCert is: ", jwtPublicCert)
}

func logErr(level, message string, err error) {
	if err != nil {
		switch level {
		case "Trace":
			logger.Trace(message, err)
		case "Debug":
			logger.Debug(message, err)
		case "Info":
			logger.Info(message, err)
		case "Warning":
			logger.Warning(message, err)
		case "Error":
			logger.Error(message, err)
		case "Fatal":
			logger.Fatal(message, err)
			os.Exit(1)
		case "Panic":
			logger.Panic(message, err)
			os.Exit(1)
		default:
			logger.Info("LogErr: Unknown Log level "+level+"( "+message+") ", err)
		}
	}
}

type myClaims struct {
	Nbf int64 `json:"nbf"`
	jwt.StandardClaims
}

// AddFooBarHeader adds custom "Foo: Bar" header to the request
func AddFooBarHeader(rw http.ResponseWriter, r *http.Request) {
	logger.Info("Processing HTTP request in Golang plugin!!")
	claims := myClaims{
		Nbf: time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(jwtPrivateKey)
	logErr("Info", "token.SignedString(jwtPrivateKey) ", err)
	/*
		t := jwt.New(jwt.GetSigningMethod("RS256"))
		t.Claims["AccessToken"] = "level1"
		t.Claims["CustomUserInfo"] = struct {
			Name string
			Kind string
		}{user, "human"}
		t.Claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
		tokenString, err := t.SignedString(signKey) */

	r.Header.Add("JWT", ss)
}

func main() {}

// docker run --rm -v /C//Users/pstubbs/go/src/mine/tyk-plugins/jwt:/plugin-source tykio/tyk-plugin-compiler:v2.9.4.2 tyk-jwt.so
// cp .\tyk-jwt.so C:\Users\pstubbs\tyk\plugins\2.9.4.2\
// docker container restart sandbox-1
