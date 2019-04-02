# Reserve addresses 

## Get all addresses 

```shell
curl -X GET "http://gateway.local/addresses"
```

> the above request will return reponse like this:

```json
[
    {
        "timestamp": 1518038157000,
        "id": 1,
        "address": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
        "type": "reserve",
        "description": "Kyber network reserve"
    },
    {
        "timestamp": 1530815752000,
        "id": 2,
        "address": "0x21433dec9cb634a23c6a4bbcce08c83f5ac2ec18",
        "type": "reserve",
        "description": "Kyber network reserve 2"
    }
]
```

### HTTP request

`GET http://gateway.local/addresses`


## Get address by id

```shell
curl -X GET "http://gateway.local/addresses/1"
```

> the above request will return reponse like this:

```json
{
    "timestamp": 1518038157000,
    "id": 1,
    "address": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
    "type": "reserve",
    "description": "Kyber network reserve"
}
```

### HTTP request

`GET http://gateway.local/addresses/:id`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
id | integer | true | none | id of reserve


## Added new address

```shell
curl -X POST "http://gateway.local/addresses" \ 
-H "Content-Type: application/json" \ 
-D '{"address": "0x63825c174ab367968EC60f061753D3bbD36A0D8F", "type":"reserve", "description": "Kyber network reserve"}' 
```

> the above request will return reponse like this:

```json
{
    "id": 1
}
```

### HTTP request

`POST http://gateway.local/addresses`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
address | integer | true | none | address value 
type | string | true | including: "reserve", "pricing operator", "sanity operator", "intermediate operator", "centralized exchange deposit addresses", 
description | string | false | empty | description of the reserve address 

## Update an address

```shell
curl -X PUT "http://gateway.local/addresses/1" \ 
-H "Content-Type: application/json" \ 
-D '{"address": "0x63825c174ab367968EC60f061753D3bbD36A0D8F", "type":"reserve", "description": "Kyber network reserve"}' 
```

> the above request will return reponse like this:

```json
{
    "id": 1
}
```

### HTTP request

`PUT http://gateway.local/addresses/:id`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
id | integer | true | none | 
address | string | true | none | address value 
type | string | true | including: "reserve", "pricing operator", "sanity operator", "intermediate operator", "centralized exchange deposit addresses", 
description | string | false | empty | description of the reserve address 