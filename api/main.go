package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var MySigningKey []byte

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	MySigningKey = []byte(os.Getenv("SECRET_KEY"))

}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super secret Information")
}

func isAuthorized(endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token,err:=jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Invalid signing method")
				}
				aud := "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

				if !checkAudience {
					return nil, fmt.Errorf("invalid audience")
				}
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf("Invalid iss")
				}
				return MySigningKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "No authorization token provided")
		}
	})
}
func handleRequests() {
	fmt.Println(MySigningKey)
	http.Handle("/", isAuthorized(homePage))

	log.Fatal(http.ListenAndServe(":9001", nil))
}
func main() {
	fmt.Println("server")

	handleRequests()
}
