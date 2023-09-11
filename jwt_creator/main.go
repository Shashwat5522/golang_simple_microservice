package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var MySigningKey []byte

func GetJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = "shashwat"
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	tokenString, err := token.SignedString(MySigningKey)
	if err != nil {
		fmt.Errorf("something went wrong:%s", err.Error())
		return "",err
	}
	return tokenString,nil
}

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	MySigningKey = []byte(os.Getenv("SECRET_KEY"))
}
func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(validToken))
}
func handleRequests() {
	http.HandleFunc("/", Index)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
func main() {
	handleRequests()
}
