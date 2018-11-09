## Trade summary

```shell
curl -X GET "http://gateway.local/trade-summary?from=1541635200000&to=1541721600000"
```

> the above request will return a struct like this

```json

```

Trade summary will summary of trade logs **from** a point of time **to** another point of time

### HTTP Request

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | start query time
to | integer | false | now | end query time

