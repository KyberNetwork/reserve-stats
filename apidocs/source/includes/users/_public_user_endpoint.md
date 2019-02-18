# Public user endpoint

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