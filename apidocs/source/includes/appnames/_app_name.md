# App names

## Create App

```shell
curl -X POST "https://gateway.local/application-names"
-H 'Content-Type: application/json'
-d '{
    "app_name": "first_app",
    "addresses": [
        "0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700",
        "0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
    ]
}'
```

> sample response

```json
{
    "id": 1,
    "app_name": "first_app",
    "addresses": [
        "0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700",
        "0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
    ]
}
```

### HTTP Request

`POST https://gateway.local/application-names`


## Update App

```shell
curl -X PUT "https://gateway.local/application-names/1"
-H 'Content-Type: application/json'
-d '{
    "app_name": "first_app",
    "addresses": [
        "0xC26633E672b6A750dC06109be7f2C5dFe36670d1"
    ]
}'
```

> sample response

```json
```

### HTTP Request

`PUT https://gateway.local/application-names/:appID`


## Get all app

```shell
curl -X GET "https://gateway.local/application-names"
```

> sample response

```json
[{
    "id": 1,
    "app_name": "first_app",
    "addresses": [
        "0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700",
        "0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
    ]
},{
    "id": 2,
    "app_name": "second_app",
    "addresses": [
        "0xC26633E672b6A750dC06109be7f2C5dFe36670d1"
    ]
}
]
```

### HTTP Request

`GET https://gateway.local/application-names`

## Get an app by id

```shell
curl -X GET "https://gateway.local/application-names/1"
```

> sample response

```json
{
    "id": 1,
    "app_name": "first_app",
    "addresses": [
        "0x3baE9b9e1dca462Ad8827f62F4A8b5b3714d7700",
        "0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
    ]
}
```

### HTTP Request

`GET https://gateway.local/:appID`


## Delete an app

```shell
curl -X DELETE "https://gateway.local/1"
```

> sample response

```json
```

### HTTP Request

`DELETE https://gateway.loal/:appID`