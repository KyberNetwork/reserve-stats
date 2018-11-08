# Users

User component provide service for user dashboard can update user kyced addresses list, then reserve stats server can store and then return appropriate statistic number

## Update user address

```shell

```

>  sample response:

> the endpoint will return http code 200 when success, otherwise it will return error

```json
{
    "error": <error>
}
```

### HTTP request:

`POST http://localhost:8006/users`

Content-type: application/json


## Get user info by address

```shell
curl -X GET "http://localhost:8002/users?address=0x829bd824b016326a401d083b33d092293333a830"
```

> the above request will return reponse like this:

```json
{
    "kyc": false,
    "cap": 12313471038132,
    "rich": false
}
```

### HTTP request

`GET http://localhost:8002/users`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
address | string | true | null | address of user who you want to check info