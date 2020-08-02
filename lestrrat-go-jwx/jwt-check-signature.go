package main

import (
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat/go-jwx/jwk"
)

const token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjE0MTYzMzY2ODc5MDA2OTQ5MDI5IiwidHlwIjoiSldUIn0.eyJhdWQiOlsiR29sYW5nIFVzZXJzIl0sImlhdCI6MjMzNDMxMjAwLCJzdWIiOiJodHRwczovL2dpdGh1Yi5jb20vbGVzdHJyYXQtZ28vand4L2p3dCIsInByaXZhdGVDbGFpbUtleSI6IkhlbGxvLCBXb3JsZCEifQ.eNno9r81dQczI4O706_i_PLyjkEE8_VYYg3yHVqcaXDTbhIJankQA_DVejOE654tlOnJohnmTJH6Nix-HEiKJ9EsExkEYn93AfnpEg0C9Mihbxheccm5p7Cf8adLTD9gWPD85CONkFnw9Hkwc9SvVkLcMuB7OTAMadNWrV78shKZg-4s-Dn_AjLi2AEBDMsNURHZyK0HTzM1QO0ZKUVrIBbaxVfE9O3b06JuA6c3EY_DLwgBQ2zUnL1rQ3o42n5o4WNb8cbqRaH_jkZ3efWBvwDtShB6EE8tcPAlSafnl2CZPD-VF5GdL3eRvVeXoY9uSqTVkyD1zyTabhtH-taK9QZfBK9Fbt-5wSuHQYZuqGZcL_MHsU06QQgvOpcgj7Zxc10GrZZDWIeUhcdsUfrZ2dAD8BaoL_6cJUZm7ywT3Ito5Rzmclx0VypiZMXsqBdyhPBXzO4gAV_QqCxXJInUj_VYowpAF4pG3qvjXFWZxOICA_nJmmXRPr9Kw6ZkvyzA06IgQflml4yy1-iMf8PUXblyEf_Q5x2ArdSUm7S9lX5PntZt5K0vAYSkNxNq4rPhefLO8j4z6T9tbQxqQ0C0aZS6T6a32D2tolZyDBKXFOyF4eBkoW3ztIVHPvm0voly5UKSs3LZWpP9XCSmo3ADSi78fDlzX1_6sbHMAtSgVxI"

const jwksURL = "http://192.168.1.81:8080"

func getKey(token *jwt.Token) (interface{}, error) {

	// TODO: cache response so we don't have to make a request every time
	// we want to verify a JWT
	set, err := jwk.FetchHTTP(jwksURL)
	fmt.Println("set")
	fmt.Println(spew.Sdump(set))
	fmt.Println("")
	if err != nil {
		return nil, err
	}
	fmt.Println("token")
	fmt.Println(spew.Sdump(token))
	fmt.Println("")

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	if key := set.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, errors.New("unable to find key")
}

func main() {
	token, err := jwt.Parse(token, getKey)
	if err != nil {
		panic(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	for key, value := range claims {
		fmt.Printf("%s\t%v\n", key, value)
	}
}
