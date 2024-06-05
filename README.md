# Mock HTTP Server

Small and simple mock HTTP server configurable with JSON.

Features:

- The server responds with configured response
- `POST` and `PUT` endpoints also have request validation and matching
  - Validation is defined by [go-playground/validator](https://github.com/go-playground/validator) tags.
  - Matching compares request JSON to configured expected JSON
- Error preparing: when you make a request to configured `prepareErrorPath`, the next request to prepared endpoint will return the prepared error

## Prepared error

Example of prepare error request:

POST

```json
{
  "method": "POST",
  "path": "/example-path",
  "status": 401,
  "response": {
    "mock": "example-mock-response",
    "info": { "error": true }
  }
}
```

## Config

The server is configured with `config.json` file.

If you are mounting a config with docker or have the file anywhere else but project root folder,
you can use ENV variable `CONFIG_PATH` so override it.

### Config example

```json
{
  "prepareErrorPath": "/prepareError",
  "handlers": [
    {
      "path": "/example-path",
      "method": "POST",
      "response": {
        "info": { "success": true },
        "data": "post validate request success data"
      },
      "request": {
        "validate": {
          "example": "required,lowercase,min=5"
        },
        "match": {
          "example": "mock-request-body"
        }
      }
    }
  ]
}
```

Serve port is configured by env `PORT`, but defaults to `:8080`

## How to run

## Docker

```sh
docker pull tmavrin/mock-http-server
```

```sh
docker run --name=mock-server -e PORT=8080 -e CONFIG_PATH='/path/to/config' tmavrin/mock-http-server
```

Available ENVIRONMENT variables:

- `PORT` -> defaults to `8080`
- `CONFIG_PATH` -> defaults to example config from the repo (`./config.json`)

## Repo clone

You can clone the repo and run the server with go (go 1.22 required)

```sh
go run .
```
