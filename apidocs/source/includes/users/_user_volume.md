## User Volume

```shell
curl -X GET "http://gateway.local/user-volume?from=1541548800000&to=1541635200000&userAddr=0x08530fEC01exxx&freq=d"
```

> sample response:

```json
{
    "1541635200000": {
        "eth_amount": 0.25758993303406674,
        "usd_amount": 56.41625875145447
    }
}
```

API return volume of an user in a time range with a frequency (daily or hourly)

### HTTP Request

`GET http://gateway.local/user-volume`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | 
to | integer | false | now |
userAddr | string | true | empty | user address to query data for