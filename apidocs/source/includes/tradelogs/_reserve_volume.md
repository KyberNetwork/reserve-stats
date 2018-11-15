## Reserve volume

```shell
curl -X GET "http://gateway.local/reserve-volume?from=1519862400000&to=1519948800000&reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F&asset=ETH&freq=d
```

> above request will return struct like this

```json
{
    "1519862400000": {
        "eth_amount": 110.47715152158673,
        "usd_amount": 93504.01774619528,
        "volume": 110.47715152158673
    },
    "1519948800000": {
        "eth_amount": 73.37498746994575,
        "usd_amount": 62785.41974102476,
        "volume": 73.37498746994575
    }
}
```

This api return volume of one reserve by **asset** amount and **usd** amount


### HTTP Request

`GET http://gateway.local/reserve-volume`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | start time to query (millisecond)
to | integer | false | now | end time to query (millisecond) 
reserve | string | true | empty | reserve address
asset | string | true | empty | asset to get volume for
freq | string | false | h (hour) | frequency to get aggregated data for (h - hour, d - day)
