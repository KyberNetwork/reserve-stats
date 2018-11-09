## Reserve volume

```shell
curl -X GET "http://gateway.local/reserve-volume?from=1541548800000&to1541635200000"
```

> above request will return struct like this


### HTTP Request

`GET http://gateway.local/reserve-volume`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | start time to query (millisecond)
to | integer | false | now | end time to query (millisecond) 

