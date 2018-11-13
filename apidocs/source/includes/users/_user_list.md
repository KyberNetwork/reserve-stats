## User list api

<aside class="notice">Signing required</aside>

```shell
curl -X GET "http://gateway.local/user-list?from=1541548800000&to=1541635200000"
```

> sample response:

```json
[
    {
        "user_address": "0x003933989385fAb554623D558523cfD23Cea8503",
        "total_eth_volume": 0.2,
        "total_usd_volume": 43.832177416861725
    },
    {
        "user_address": "0x00EC493b759760685c18870F3FDB39be6f51C017",
        "total_eth_volume": 1.0045553389337065,
        "total_usd_volume": 220.15923920598942
    },
    {
        "user_address": "0x028ab225A8224475b84c24535b7325fbD56f0390",
        "total_eth_volume": 1.99,
        "total_usd_volume": 436.13016529777417
    }
]
```

Return list of user who trade in a time range and their volume.

### HTTP Request

`GET http://gateway.local/user-list`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | time to query from
to | integer | false | now | time to query to