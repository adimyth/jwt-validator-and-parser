package main

import (
	"fmt"
	"strings"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/golang-jwt/jwt"
)

type Config struct {
	JWTSecret string   `json:"jwt_secret"`
	UserKeys  []string `json:"user_keys"`
}

func New() interface{} {
	return &Config{UserKeys: []string{"first_name", "last_name", "role_code"}}
}

func (config Config) Access(kong *pdk.PDK) {
	// Read Authorization Bearer token
	tokenString, err := kong.Request.GetHeader("Authorization")
	if err != nil {
		kong.Log.Info("could not find `Bearer` token: %s", err)
	}

	// Remove `Bearer`
	tokenString = strings.Split(tokenString, "Bearer ")[1]

	// Signature validation & decoding claims from token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			kong.Log.Info("unable to parse JWT")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})

	if err == nil {
		// Add request headers only if validation is successful
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			for _, k := range config.UserKeys {
				h := fmt.Sprintf("X-AUTH-%s", k)
				v := fmt.Sprintf("%v", claims["user"].(map[string]interface{})[k])
				_ = kong.ServiceRequest.SetHeader(h, v)
			}
		}
	} else {
		kong.Log.Info("error parsing JWT: ", err)
	}

}

func main() {
	Version := "0.1"
	Priority := 1
	_ = server.StartServer(New, Version, Priority)
}
