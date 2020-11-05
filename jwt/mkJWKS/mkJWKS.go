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

func translateSignatureAlgorithm(SigAlg string) (string) {
  if SigAlg == "SHA256-RSA" {
    return "RS256"
  } else {
    fmt.Println("[Fatal]Unknown SignatureAlgorithm ", SigAlg)
    os.Exit(1)
  }
  return ""
}

func main() {
	var jwks jose.JSONWebKeySet
	for _, certFile := range os.Args[1:] {
		fmt.Println("Loading " + certFile)
		certBytes, err := ioutil.ReadFile(certFile)
		if err != nil {
			fmt.Println("[FATAL]Unable to load "+certFile+": ", err)
			os.Exit(1)
		}
		var certs []*x509.Certificate
		var cert *x509.Certificate
		var block *pem.Block
		// read all the blocks from the file so assuming that they make a chain
		// the first one will control the kid.
		for len(certBytes) > 0 {
			block, certBytes = pem.Decode(certBytes)
			cert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				fmt.Println("[FATAL]Cannot parse "+certFile+", error: ", err)
				os.Exit(1)
			}
			//fmt.Println(cert.SerialNumber.String())
			certs = append(certs, cert)
		}
		cert = certs[0]
		x5tSHA1 := sha1.Sum(cert.Raw)
		x5tSHA256 := sha256.Sum256(cert.Raw)

		jwk := jose.JSONWebKey{
			Key:                         cert.PublicKey,
			KeyID:                       cert.SerialNumber.String(),
			Algorithm:                   translateSignatureAlgorithm(cert.SignatureAlgorithm.String()),
			Certificates:                certs,
			CertificateThumbprintSHA1:   x5tSHA1[:],
			CertificateThumbprintSHA256: x5tSHA256[:],
			Use:                         "sig",
		}
		jwks.Keys = append(jwks.Keys, jwk)
	}
	jsonJwks, err := json.Marshal(&jwks)
	if err != nil {
		fmt.Println("[FATAL]Unable to marshal the json: ", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonJwks))
}
