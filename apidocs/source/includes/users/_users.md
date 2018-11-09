# Users

User component provide service for user dashboard can update user kyced addresses list, then reserve stats server can store and then return appropriate statistic number

## Update user address

```shell
curl -X POST "http://gateway.local/users"
-H 'Content-Type: application/json'
-d '{
    "email": "sample@gmail.com",
    "user_info": [
        {
            "address": "0x829bd824b016326a401d083b33d092293333a830",
            "timestamp": 1538020234000
        },
        {
            "address": "0x5270e589f57936123e649c1c3aa07fd54f2fcb64",
            "timestamp": 1538975051000
        }
    ]
}'
```

>  sample response:

> the endpoint will return http code 200 when success, otherwise it will return error

```json
{
    "error": <error>
}
```

### HTTP request:

`POST http://gateway.local/users`

Content-type: application/json


## Get user info by address

```shell
curl -X GET "http://gateway.local/users?address=0x829bd824b016326a401d083b33d092293333a830"
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

`GET http://gateway.local/users`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
address | string | true | null | address of user who you want to check info