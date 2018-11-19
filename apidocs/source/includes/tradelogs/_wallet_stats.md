## Wallet Stats API

```shell
curl -X GET "http://gateway.local/wallet-stats?from=1541548800000&to=1541635200000&walletAddr=0xf1aaxxxxx
```

> sample response:

```json
```

### HTTP Request

`GET http://gateway.local/wallet-stats`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | time to query data from
to | integer | false | now | time to query data to
walletAddr | string | true | empty | wallet address to query stat for