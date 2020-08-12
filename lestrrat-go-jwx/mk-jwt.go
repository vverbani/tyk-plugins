package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

//const aLongLongTimeAgo = 233431200

var (
	errKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	errNotRSAPrivateKey    = errors.New("Key is not a valid RSA private key")
	errNotRSAPublicKey     = errors.New("Key is not a valid RSA public key")
)

func parseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		fmt.Println("ErrKeyMustBePEMEncoded", errKeyMustBePEMEncoded)
		return nil, errKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, errNotRSAPrivateKey
	}

	return pkey, nil
}

func parseRSAPrivateKeyFromFile(rsaPrivateKeyLocation string) (*rsa.PrivateKey, error) {
	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		fmt.Println("No RSA private key found: ", err)
		return nil, err
	}
	return parseRSAPrivateKeyFromPEM(priv)
}

func parseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, errNotRSAPublicKey
	}

	return pkey, nil
}

func parseRSAPublicKeyFromFile(rsaPublicKeyLocation string) (*rsa.PublicKey, error) {
	pub, err := ioutil.ReadFile(rsaPublicKeyLocation)
	if err != nil {
		fmt.Println("No RSA public key found: ", err)
		os.Exit(1)
	}
	return parseRSAPublicKeyFromPEM(pub)
}

func parseRSACertFromPEM(key []byte) (*x509.Certificate, error) {
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errKeyMustBePEMEncoded
	}

	// Parse the cert from the PEM block
	if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
		return cert, nil
	} else {
		return nil, err
	}
}

func parseRSACertFromFile(rsaCertKeyLocation string) (*x509.Certificate, error) {
	cert, err := ioutil.ReadFile(rsaCertKeyLocation)
	if err != nil {
		fmt.Println("No RSA private key found: ", err)
		return nil, err
	}
	return parseRSACertFromPEM(cert)
}

func parseJSONFromFIle(claimsFile string) (map[string]interface{}, error) {
	jsonClaimsFile, err := os.Open(claimsFile)
	if err != nil {
		return nil, err
	}
	defer jsonClaimsFile.Close()
	jsonByteValue, _ := ioutil.ReadAll(jsonClaimsFile)
	var jsonClaims map[string]interface{}
	json.Unmarshal([]byte(jsonByteValue), &jsonClaims)
	return jsonClaims, nil
}

func createJwt(certFile, keyFile, claimsFile string) {
	cert, err := parseRSACertFromFile(certFile)
	json, err := parseJSONFromFIle(claimsFile)

	hdrs := jws.NewHeaders()
	hdrs.Set(jws.KeyIDKey, cert.SerialNumber.String())

	s := jwt.New()
	s.Set(jwt.SubjectKey, `https://github.com/lestrrat-go/jwx/jwt`)
	s.Set(jwt.AudienceKey, `Golang Users`)
	s.Set(jwt.IssuedAtKey, time.Now().Unix)
	for jsonKey, jsonValue := range json {
		s.Set(jsonKey, jsonValue)
	}

	privkey, err := parseRSAPrivateKeyFromFile(keyFile)
	if err != nil {
		log.Printf("Failed to load private key from %s: %s", keyFile, err)
		return
	}

	signed, err := jwt.Sign(s, jwa.RS256, privkey, jwt.WithHeaders(hdrs))
	if err != nil {
		log.Printf("Failed to created JWS message: %s", err)
		return
	}
	pubkey := cert.PublicKey.(*rsa.PublicKey)

	fmt.Println("Signed jws with certificate in ", certFile)
	fmt.Println(string(signed))
	fmt.Println("")

	token, err := jwt.Parse(bytes.NewReader(signed), jwt.WithVerify(jwa.RS256, pubkey))
	if err != nil {
		panic(err)
	}
	//fmt.Println(token)
	//claims := token.PrivateClaims.(jwt.MapClaims)
	fmt.Println("Private claims:")
	for key, value := range token.PrivateClaims() {
		fmt.Printf("%s\t->\t%v\n", key, value)
	}
	fmt.Println("All claims:")
	// seriously??? This is what you have to do to get a list of claims?
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for iter := token.Iterate(ctx); iter.Next(ctx); {
		pair := iter.Pair()
		fmt.Printf("%s -> %v\n", pair.Key, pair.Value)
	}

	// When you received a JWS message, you can verify the signature
	// and grab the payload sent in the message in one go:
	verified, err := jws.Verify(signed, jwa.RS256, *pubkey)
	if err != nil {
		log.Printf("Failed to verify message: %s", err)
		return
	}
	fmt.Printf("\nSigned message verified! -> %s\n", verified)
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	cert := flag.String("cert", "cert.pem", "The x509 RSA public certificate")
	key := flag.String("key", "key.pem", "The RSA private key")
	claims := flag.String("claims", "claims.json", "A file of claims in json format")
	//jwks := flag.String("jwks", "http://localhost:8080/jwks.json", "The matching JWKS to the cert and key")
	flag.Parse()
	if *cert == "" || *key == "" {
		fmt.Println("Must provide --cert, --key, --claims")
		os.Exit(1)
	}
	createJwt(*cert, *key, *claims)
}
