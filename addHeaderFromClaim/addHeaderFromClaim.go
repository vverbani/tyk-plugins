package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/TykTechnologies/tyk/headers"
	jwt "github.com/dgrijalva/jwt-go"
)

func stripBearer(token string) string {
	if len(token) > 6 && strings.ToUpper(token[0:7]) == "BEARER " {
		return token[7:]
	}
	return token
}

var policyArray = []string{"60d30bc4077b5f00175ad141", "60d30bc4077b5f00175ad142", "60d30bc4077b5f00175ad143", "60d30bc4077b5f00175ad144"}
var headerVal = []string{"Val1", "Val2", "Val3", "Val4"}
var headerField = "X-Tyk-Pol-Val"

// AddHeaderFromClaim based on Policy ID
func AddHeaderFromClaim(rw http.ResponseWriter, r *http.Request) {

	// 1. Retrieve the 'pol' claim from the JWT used to authenticate
	// 2. Look that up in policyArray
  // 3. Find the corresponding headerVal into the the header headerField

	log.Println("Start addHeaderFromClaim Plugin")

	// decode the JWT from the Authorization header
	jwtStr := stripBearer(r.Header.Get(headers.Authorization))
	token, _ := jwt.Parse(jwtStr, nil)
	claims, _ := token.Claims.(jwt.MapClaims)

	// check if the pol claim is populated and is in the array policyArray
	if val, ok := claims["pol"]; ok {
		polStr := fmt.Sprintf("%v", val)
		log.Println("Pol is " + polStr)
		for k, v := range policyArray {
			if v == polStr {
				// add the headerField header corresponding to the policy
				r.Header.Add(headerField, headerVal[k])
				break
			}
		}
	}

	log.Println("End addHeaderFromClaim Plugin")

}

func main() {}
