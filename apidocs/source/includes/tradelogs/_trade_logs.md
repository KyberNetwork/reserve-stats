
# Tradelogs 

Trade logs component crawls all tradelogs which go through Kyber contract in order to store and aggregation then provide appropriate stats.

## Get trade logs 

```shell
curl -X GET "http://gateway.local/trade-logs
```

> The above command returns JSON structured like this:

```json
[
    {
        "timestamp": 1541730585000,
        "block_number": 6669829,
        "tx_hash": "0x207d546e61ac1bace8ca1b8cd0075bbcc300437f1dfd84d6cb1a91d356aee915",
        "eth_receival_sender": "0x0000000000000000000000000000000000000000",
        "eth_receival_amount": 0,
        "user_addr": "0xc053dd595f5c6c4660c49091486e643e64fcd404",
        "src_addr": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
        "dst_addr": "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359",
        "src_amount": 100000000000000000,
        "dst_amount": 21246750000000000000,
        "fiat_amount": 21.340763439234546,
        "burn_fees": [
            {
                "reserve_addr": "0x9e2b650f890236ab49609c5a6b00cddb4e61f408",
                "amount": 62500000000000000
            }
        ],
        "wallet_fees": [
            {
                "reserve_addr": "0x9e2b650f890236ab49609c5a6b00cddb4e61f408",
                "wallet_addr": "0xb9e29984fe50602e7a619662ebed4f90d93824c7",
                "amount": 62500000000000000
            }
        ],
        "ip": "",
        "country": ""
    },
    {
        "timestamp": 1541731218000,
        "block_number": 6669864,
        "tx_hash": "0xff2f34548244060cff2aa6c2a244f76b80617409b9d08d20517c74e4d1cc20eb",
        "eth_receival_sender": "0x0000000000000000000000000000000000000000",
        "eth_receival_amount": 0,
        "user_addr": "0xe25f97ba26ec8be13f54da984da151366bb60bd0",
        "src_addr": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
        "dst_addr": "0x4f3afec4e5a3f2a6a1a411def7d7dfe50ee057bf",
        "src_amount": 14110000000000000000,
        "dst_amount": 73069182000,
        "fiat_amount": 3011.1817212759943,
        "burn_fees": [
            {
                "reserve_addr": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
                "amount": 17637500000000000000
            }
        ],
        "wallet_fees": null,
        "ip": "186.109.142.254",
        "country": "AR"
    }
]
```

Return list of trade logs **from** a point time and **to** another point of time

### HTTP Request

`GET http://gateway.local/trade-logs`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | start time to query trade logs
to | integer | false | now | end time to query trade logs
