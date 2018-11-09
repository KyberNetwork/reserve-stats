## Wallet fee

```shell
curl -X GET "http://gateway.local/wallet-fee?from=1530921600000&to=1531008000000&freq=d&reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F&walletAddr=0xf1aa99c69715f423086008eb9d06dc1e35cc504d"
```

> sample response 

```json
{
    "1541548800000": 11.28885889129255,
    "1541635200000": 24.450580410254613
}
```

### HTTP Request

`GET http://gateway.local/wallet-fee`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | 
to | integer | false | now | 