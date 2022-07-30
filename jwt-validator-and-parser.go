package main

import (
	"fmt"
	"strings"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/golang-jwt/jwt"
)

type Config struct {
	UserKeys []string `json:"user_keys"`
}

func New() interface{} {
	return &Config{UserKeys: []string{"first_name", "last_name", "role_code"}}
}

// TODO: Verify auth token
func (config Config) Access(kong *pdk.PDK) {
	// Read Authorization Bearer token
	tokenString, err := kong.Request.GetHeader("Authorization")
	if err != nil {
		kong.Log.Err("could not find `Bearer` token: %s", err)
	}

	// Remove `Bearer`
	tokenString = strings.Split(tokenString, "Bearer ")[1]

	// Decode claims from token
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		kong.Log.Err("error decoding claims from token: ", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		for _, k := range config.UserKeys {
			v := fmt.Sprintf("%v", claims["user"].(map[string]interface{})[k])
			h := fmt.Sprintf("X-AUTH-%s", k)
			_ = kong.Response.SetHeader(h, v)
		}
	}
}

func main() {
	Version := "0.1"
	Priority := 1
	_ = server.StartServer(New, Version, Priority)
}
