package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
)

var (
	errKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	errNotRSAPrivateKey    = errors.New("Key is not a valid RSA private key")
	errNotRSAPublicKey     = errors.New("Key is not a valid RSA public key")
)

func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		fmt.Println("ErrKeyMustBePEMEncoded", errKeyMustBePEMEncoded)
		return nil, errKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		fmt.Println("ParsePKCS1PrivateKey")
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			fmt.Println("ParsePKCS8PrivateKey")
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		fmt.Println("ErrNotRSAPrivateKey", errNotRSAPrivateKey)
		return nil, errNotRSAPrivateKey
	}

	return pkey, nil
}

func ParseRSAPrivateKeyFromFile(rsaPrivateKeyLocation string) (*rsa.PrivateKey, error) {
	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		fmt.Println("No RSA private key found: ", err)
		os.Exit(1)
	}
	return ParseRSAPrivateKeyFromPEM(priv)
}

func ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		fmt.Println("ParsePKIXPublicKey")
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			fmt.Println("ParseCertificate")
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

func ParseRSAPublicKeyFromFile(rsaPublicKeyLocation string) (*rsa.PublicKey, error) {
	pub, err := ioutil.ReadFile(rsaPublicKeyLocation)
	if err != nil {
		fmt.Println("No RSA public key found: ", err)
		os.Exit(1)
	}
	return ParseRSAPublicKeyFromPEM(pub)
}

func ParseRSACertFromPEM(key []byte) (*x509.Certificate, error) {
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errKeyMustBePEMEncoded
	}

	// Parse the cert
	if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
		return cert, nil
	} else {
		return nil, err
	}
}

func ParseRSACertFromFile(rsaCertKeyLocation string) (*x509.Certificate, error) {
	cert, err := ioutil.ReadFile(rsaCertKeyLocation)
	if err != nil {
		fmt.Println("No RSA private key found: ", err)
		os.Exit(1)
	}
	return ParseRSACertFromPEM(cert)
}

func Example_jwt() {
	const aLongLongTimeAgo = 233431200

	cert, err := ParseRSACertFromFile("cert.pem")
	kid := cert.SerialNumber.String()
	hdrs := jws.NewHeaders()
	hdrs.Set(jws.KeyIDKey, kid)

	t := jwt.New()
	t.Set(jwt.SubjectKey, `https://github.com/lestrrat-go/jwx/jwt`)
	t.Set(jwt.AudienceKey, `Golang Users`)
	t.Set(jwt.IssuedAtKey, time.Unix(aLongLongTimeAgo, 0))
	t.Set(`privateClaimKey`, `Hello, World!`)

	buf, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		fmt.Printf("failed to generate JSON: %s\n", err)
		return
	}

	fmt.Printf("%s\n", buf)
	fmt.Printf("aud -> '%s'\n", t.Audience())
	fmt.Printf("iat -> '%s'\n", t.IssuedAt().Format(time.RFC3339))
	if v, ok := t.Get(`privateClaimKey`); ok {
		fmt.Printf("privateClaimKey -> '%s'\n", v)
	}
	fmt.Printf("sub -> '%s'\n", t.Subject())

	//key, err := rsa.GenerateKey(rand.Reader, 2048)
	key, err := ParseRSAPrivateKeyFromFile("key1.pem")

	//fmt.Println("key")
	//fmt.Println(spew.Sdump(key))
	if err != nil {
		fmt.Printf("failed to generate private key: %s", err)
		return
	}

	{
		// Signing a token (using raw rsa.PrivateKey)
		signed, err := jwt.Sign(t, jwa.RS256, key, jwt.WithHeaders(hdrs))
		if err != nil {
			fmt.Printf("failed to sign token: %s", err)
			return
		}
		_ = signed
		fmt.Println("Signing a token (using raw rsa.PrivateKey)")
		fmt.Println(string(signed))
		fmt.Println("")
	}

	{
		// Signing a token (using JWK)
		jwkKey, err := jwk.New(key)
		if err != nil {
			fmt.Printf("failed to create JWK key: %s", err)
			return
		}

		signed, err := jwt.Sign(t, jwa.RS256, jwkKey, jwt.WithHeaders(hdrs))
		if err != nil {
			fmt.Printf("failed to sign token: %s", err)
			return
		}
		_ = signed
		fmt.Println("Signing a token (using using JWK)")
		fmt.Println(string(signed))
		fmt.Println("")

	}
	// OUTPUT:
	// {
	//   "aud": [
	//     "Golang Users"
	//   ],
	//   "iat": 233431200,
	//   "sub": "https://github.com/lestrrat-go/jwx/jwt",
	//   "privateClaimKey": "Hello, World!"
	// }
	// aud -> '[Golang Users]'
	// iat -> '1977-05-25T18:00:00Z'
	// privateClaimKey -> 'Hello, World!'
	// sub -> 'https://github.com/lestrrat-go/jwx/jwt'
}

func Example_openid() {
	const aLongLongTimeAgo = 233431200

	t := openid.New()
	t.Set(jwt.SubjectKey, `https://github.com/lestrrat-go/jwx/jwt`)
	t.Set(jwt.AudienceKey, `Golang Users`)
	t.Set(jwt.IssuedAtKey, time.Unix(aLongLongTimeAgo, 0))
	t.Set(`privateClaimKey`, `Hello, World!`)

	addr := openid.NewAddress()
	addr.Set(openid.AddressPostalCodeKey, `105-0011`)
	addr.Set(openid.AddressCountryKey, `日本`)
	addr.Set(openid.AddressRegionKey, `東京都`)
	addr.Set(openid.AddressLocalityKey, `港区`)
	addr.Set(openid.AddressStreetAddressKey, `芝公園 4-2-8`)
	t.Set(openid.AddressKey, addr)

	buf, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		fmt.Printf("failed to generate JSON: %s\n", err)
		return
	}
	fmt.Printf("%s\n", buf)

	t2, err := jwt.ParseBytes(buf, jwt.WithOpenIDClaims())
	if err != nil {
		fmt.Printf("failed to parse JSON: %s\n", err)
		return
	}
	if _, ok := t2.(openid.Token); !ok {
		fmt.Printf("using jwt.WithOpenIDClaims() creates an openid.Token instance")
		return
	}
}

func Example_jwk() {
	set, err := jwk.FetchHTTP("https://foobar.domain/jwk.json")
	if err != nil {
		fmt.Printf("failed to parse JWK: %s", err)
		return
	}

	// If you KNOW you have exactly one key, you can just
	// use set.Keys[0]
	keys := set.LookupKeyID("mykey")
	if len(keys) == 0 {
		fmt.Printf("failed to lookup key: %s", err)
		return
	}

	var key interface{} // This is the raw key, like *rsa.PrivateKey or *ecdsa.PrivateKey
	if err := keys[0].Raw(&key); err != nil {
		fmt.Printf("failed to create public key: %s", err)
		return
	}

	// Use key for jws.Verify() or whatever
	_ = key
}

func Example_jws() {
	//privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	privkey, err := ParseRSAPrivateKeyFromFile("key1.pem")
	if err != nil {
		fmt.Printf("failed to generate private key: %s", err)
		return
	}

	buf, err := jws.Sign([]byte("Lorem ipsum"), jwa.RS256, privkey)
	if err != nil {
		fmt.Printf("failed to created JWS message: %s", err)
		return
	}
	fmt.Println(string(buf))

	// When you received a JWS message, you can verify the signature
	// and grab the payload sent in the message in one go:
	verified, err := jws.Verify(buf, jwa.RS256, &privkey.PublicKey)
	if err != nil {
		fmt.Printf("failed to verify message: %s", err)
		return
	}

	fmt.Printf("signed message verified! -> %s", verified)
}

func Example_jwe() {
	//privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	privkey, err := ParseRSAPrivateKeyFromFile("key1.pem")
	if err != nil {
		fmt.Printf("failed to generate private key: %s", err)
		return
	}
	//fmt.Println(&privkey.PublicKey)

	payload := []byte("Lorem Ipsum")

	encrypted, err := jwe.Encrypt(payload, jwa.RSA1_5, &privkey.PublicKey, jwa.A128CBC_HS256, jwa.NoCompress)
	if err != nil {
		fmt.Printf("failed to encrypt payload: %s", err)
		return
	}
	fmt.Println(string(encrypted))

	decrypted, err := jwe.Decrypt(encrypted, jwa.RSA1_5, privkey)
	if err != nil {
		fmt.Printf("failed to decrypt: %s", err)
		return
	}
	fmt.Println(string(decrypted))

	if string(decrypted) != "Lorem Ipsum" {
		fmt.Printf("WHAT?!")
		return
	}
}

func main() {
	fmt.Println("Start Example_jwe")
	Example_jwe()
	fmt.Println("")

	fmt.Println("Example_jws")
	Example_jws()
	fmt.Println("")

	fmt.Println("Example_jwt")
	Example_jwt()
	fmt.Println("")
}
