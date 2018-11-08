---
title: API Reference

language_tabs: # must be one of https://git.io/vQNgJ
  - shell

toc_footers:
  - <a href='#'>Sign Up for a Developer Key</a>
  - <a href='https://github.com/lord/slate'>Documentation Powered by Slate</a>

includes:
  - errors

search: true
---

# Introduction


# Authentication

All APIs that are marked with (signing required) must follow authentication mechanism below:

1. Must be urlencoded (**x-www-form-urlencoded**)
1. Must have `signed` header with value equals to `hmac512(secret, message)`
1. Must contain `nonce` param, its value is the unix time in millisecond, it must not be before or after server time by 10s
1. `message` is constructed in following way: all query params (nonce is included) and body key-values are merged into one urlencoded string with keys are sorted.
1. `secret` is configured secret string.

Example:  

- param query: `amount=0xde0b6b3a7640000&nonce=1514554594528&token=KNC`  
- secret: `vtHpz1l0kxLyGc4R1qJBkFlQre5352xGJU9h8UQTwUTz5p6VrxcEslF4KnDI21s1`  
- signed string: `2969826a713d13b399dd0d016dad3e95949aa81ed8703ec0258abebb5f0288b96272eef68275f12a32f7e396de3b5fd63ed12b530385e08e1b676c695aacb93b`

# General 

## Get time server 


```shell
curl -X GET "http://localhost:8000/timeserver"
```

> The above command returns JSON structured like this:

```json
{
  "data": "1517479497447",
  "success": true
}
```

This endpoint return current time server

### HTTP Request

`GET http://example.com/timeserver`

