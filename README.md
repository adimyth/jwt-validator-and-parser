# JWT Validator & Parser

A kong plugin for JWT validation & parsing

![Kong JWT Validator & Parser Plugin](https://imgur.com/z7Nmuh9.png)

## JWT Validator

Validates the JWT token given _"JWT Secret"_ in plugin configuration.
Does nothing if the token is invalid. Parses & attaches _"User Keys"_ in headers & passes it onto the downstream services.

> Validates if the token has invalid signature or is expired

---

## JWT Parser

Parses claims from `Authorization: Bearer` & adds them to the Header

### Configuration

You can add the plugin with the following request:

```bash
$ curl -X POST http://kong:8000/apis/{api}/plugins \
    --data "name=jwt-validator-and-parser" \
    --data "config.user_keys=first_name,last_name,role_code"
    --data "config.jwt_secret=JWT_SECRET"
```

| Configuration | Default                        | Possible values           | Description                       |
| ------------- | ------------------------------ | ------------------------- | --------------------------------- |
| `user_keys`   | first_name,last_name,role_code | user claims in auth token | User claims                       |
| `jwt_secret`  |                                |                           | JWT secret used to sign the token |

## Deployment

### Production

```bash
docker compose up
```

### Development Mode

If you wish to make changes to the plugin, you will need to rebuild the application. Make desired changes in `jwt-validator-and-parser.go`, then rebuild & run -

```bash
rm jwt-validator-and-parser && GOOS=linux GOARCH=amd64 go build . && docker compose up
```

## Demo
#### Create a simple server
To demonstrate, create a simple server that listens on port 3000. Add a simple `GET /` endpoint & log the request headers. Example -

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		headersString, _ := json.MarshalIndent(headers, "", "  ")
		fmt.Print(string(headersString))
		fmt.Print("\n\n")
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
```
#### Make request
Send a request with Bearer token to this endpoint.

1️⃣ *Valid Token*

If a valid token is provided, the plugin adds `X-AUTH-*` tokens. In the below example, there are 5 additional headers as per the `user_keys` set in the configuration

```bash
{
"Accept": "*/*",
"Accept-Encoding": "gzip, deflate, br",
"Authorization": "Bearer eyJhbG....Om51bQ",
"Connection": "keep-alive",
"Host": "docker.for.mac.host.internal:3000",
"Postman-Token": "9cf6c76d-5a52-40a7-8396-9bad38c049d3",
"User-Agent": "PostmanRuntime/7.29.2",
"X-Auth-Email": "aditya@fnp.dev",
"X-Auth-First_name": "Aditya",
"X-Auth-Last_name": "Mishra",
"X-Auth-Phone": "9029080380",
"X-Auth-Role_code": "[company-admin appuser]",
"X-Forwarded-For": "172.29.0.1",
"X-Forwarded-Host": "localhost",
"X-Forwarded-Path": "/echo/",
"X-Forwarded-Port": "8000",
"X-Forwarded-Prefix": "/echo",
"X-Forwarded-Proto": "http",
"X-Real-Ip": "172.29.0.1"
}
```

2️⃣ *Invalid / Expired Token*

Notice that in case of invalid or expired tokens, the plugin does not add any `X-AUTH-*` headers

```bash
{
"Accept": "*/*",
"Accept-Encoding": "gzip, deflate, br",
"Authorization": "Bearer eyJhbG.....dwCt7g",
"Connection": "keep-alive",
"Host": "docker.for.mac.host.internal:3000",
"Postman-Token": "711c01d5-e9f4-4643-9fc1-032eee6b8aac",
"User-Agent": "PostmanRuntime/7.29.2",
"X-Forwarded-For": "172.29.0.1",
"X-Forwarded-Host": "localhost",
"X-Forwarded-Path": "/echo/",
"X-Forwarded-Port": "8000",
"X-Forwarded-Prefix": "/echo",
"X-Forwarded-Proto": "http",
"X-Real-Ip": "172.29.0.1"
}
```
