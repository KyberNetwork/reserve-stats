# Kyber IPLocator

## Build

```shell
    cd cmd
    go build - v .
```

## Run

```
    cd cmd
    ./cmd --data-dir=.
```

### Available flags

**data-dir**: where to load/store db file

## Available API

### Get user info by address

```http
GET http://localhost:8001/ip/:ip
```

response: Return country code

**Sample response:**

```json
{
    "country": "US"
}
```