# Reserve addresses 

Returns Ethereum addresses related to KyberNetwork's reserve. An address will have one of following types:

- **reserve**: address of KyberNetwork's reserve contract, example: 0x63825c174ab367968EC60f061753D3bbD36A0D8F
- **pricing operator**: operator addresses of conversion rate contract. Can be get by calling getOperators method of conversion rate contract, example:
0x8bc3da587def887b5c822105729ee1d6af05a5ca, 0x9224016462b204c57eb70e1d69652f60bcaf53a8
- **sanity operator**: operator address of sanity rates contract. Can be get by calling getOperators method of sanity rates
contract
- **intermediate operator**: some centralized exchanges (Huobi) does not allow deposit/withdraw directly to contract
account. Therefore, we need to use a intermediate account when deposit/withdraw funds to reserve contract.
- **deposit operator**: operator responsible for doing deposit directly into centralized exchanges
- **centralized exchange deposit addresses:** Ethereum address to deposit funds to centralized exchanges
(binance: 0x44d34a119ba21a42167ff8b77a88f0fc7bb2db90, huobi: 0x0c8fd73eaf6089ef1b91231d0a07d0d2ca2b9d66)
- **company wallet**: ethereum address for company wallet

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
-d '{"address": "0x63825c174ab367968EC60f061753D3bbD36A0D8F", "type":"reserve", "description": "Kyber network reserve"}' 
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
type | string | true | including: "reserve", "pricing_operator", "sanity_operator", "intermediate_operator", "cex_deposit_address", "company_wallet", "deposit_operator"
description | string | false | empty | description of the reserve address 

## Update an address

```shell
curl -X PUT "http://gateway.local/addresses/1" \ 
-H "Content-Type: application/json" \ 
-d '{"address": "0x63825c174ab367968EC60f061753D3bbD36A0D8F", "type":"reserve", "description": "Kyber network reserve"}' 
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
type | string | true | including: "reserve", "pricing_operator", "sanity_operator", "intermediate_operator", "cex_deposit_address", "company_wallet", "deposit_operator" 
description | string | false | empty | description of the reserve address 