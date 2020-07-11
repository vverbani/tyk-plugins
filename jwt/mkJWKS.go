package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/square/go-jose.v2"
)

func main() {
	var jwks jose.JSONWebKeySet
	for _, certFile := range os.Args[1:] {
		fmt.Println("Loading " + certFile)
		certBytes, err := ioutil.ReadFile(certFile)
		if err != nil {
			fmt.Println("[FATAL]Unable to load "+certFile+": ", err)
			os.Exit(1)
		}
		//https://gist.github.com/ukautz/cd118e298bbd8f0a88fc for multi certs in a file
		pubPEM, _ := pem.Decode(certBytes)
		if err != nil {
			fmt.Println("[FATAL]Unable to parse contents of "+certFile+"as a PEM format certificate: ", err)
			os.Exit(1)
		}

		certs, _ := x509.ParseCertificates(pubPEM.Bytes)
		for _, cert := range certs {
			x5tSHA1 := sha1.Sum(cert.Raw)
			x5tSHA256 := sha256.Sum256(cert.Raw)
			fmt.Println(cert.SerialNumber.String())
			jwk := jose.JSONWebKey{
				Key:                         cert.PublicKey,
				KeyID:                       cert.SerialNumber.String(),
				Algorithm:                   cert.SignatureAlgorithm.String(),
				Certificates:                certs,
				CertificateThumbprintSHA1:   x5tSHA1[:],
				CertificateThumbprintSHA256: x5tSHA256[:],
				Use:                         "sig",
			}
			jwks.Keys = append(jwks.Keys, jwk)
		}
	}
	jsonJwks, _ := json.Marshal(&jwks)
	fmt.Println(string(jsonJwks))
}
