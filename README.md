# JWT Validator & Parser

A kong plugin for JWT validation & parsing

## JWT Validator

⏳ Coming soon

---

## JWT Parser
Parses claims from `Authorization: Bearer` & adds them to the Header

### Configuration
You can add the plugin with the following request:

```bash
$ curl -X POST http://kong:8000/apis/{api}/plugins \
    --data "name=jwt-validator-and-parser" \
    --data "config.user_keys=first_name,last_name,role_code"
```

#### `user_keys`
Set of user keys that must be appended to the request header

`default` - first_name,last_name,role_code

You can choose from the set of user claims. Example - `email, phone, is_active`, etc

## Run locally
### Production
```bash
docker compose up
```

### Development Mode
If you wish to make changes to the plugin, you will need to rebuild the application. Make desired changes in `jwt-validator-and-parser.go`, then rebuild & run -

```bash
rm jwt-validator-and-parser && GOOS=linux GOARCH=amd64 go build . && docker compose up
```