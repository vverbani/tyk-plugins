package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

const aLongLongTimeAgo = 233431200

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

func createJwt() {
	cert, err := parseRSACertFromFile("cert.pem")
	hdrs := jws.NewHeaders()
	hdrs.Set(jws.KeyIDKey, cert.SerialNumber.String())

	s := jwt.New()
	s.Set(jwt.SubjectKey, `https://github.com/lestrrat-go/jwx/jwt`)
	s.Set(jwt.AudienceKey, `Golang Users`)
	s.Set(jwt.IssuedAtKey, time.Unix(aLongLongTimeAgo, 0))
	s.Set(`privateClaimKey`, `Hello, World!`)

	privkey, err := parseRSAPrivateKeyFromFile("key1.pem")
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return
	}

	signed, err := jwt.Sign(s, jwa.RS256, privkey, jwt.WithHeaders(hdrs))
	if err != nil {
		log.Printf("failed to created JWS message: %s", err)
		return
	}

	fmt.Println("Signed jws")
	fmt.Println(string(signed))
	fmt.Println("")

	// When you received a JWS message, you can verify the signature
	// and grab the payload sent in the message in one go:
	verified, err := jws.Verify(signed, jwa.RS256, &privkey.PublicKey)
	if err != nil {
		log.Printf("failed to verify message: %s", err)
		return
	}

	log.Printf("signed message verified! -> %s", verified)
}

func main() {
	createJwt()
}
