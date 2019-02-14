
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
    "timestamp": 1546622719000,
    "block_number": 7010000,
    "tx_hash": "0x2bd15007d29cc04495c197e9dabb4f966513013e3aa5fb86e116268def97bc34",
    "eth_amount": 2321664882262579700,
    "user_addr": "0x093ad56857aaa28615f6efef5518e564bba39c64",
    "src_addr": "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359",
    "dst_addr": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
    "src_reserve_addr": "0x44aef3101432a64d1aa16388f4b9b352b09f42a9",
    "dst_reserve_addr": "0x0000000000000000000000000000000000000000",
    "src_amount": 350000000000000000000,
    "dst_amount": 2321664882262579700,
    "fiat_amount": 340.6592798595379,
    "wallet_addr": "0xea1a7de54a427342c8820185867cf49fc2f95d43",
    "src_burn_amount": 3.8597678667615387,
    "dst_burn_amount": 0,
    "src_wallet_fee_amount": 1.654186228612088,
    "dst_wallet_fee_amount": 0,
    "integration_app": "KyberSwap",
    "ip": "",
    "country": "",
    "user_name": "",
    "profile_id": 0,
    "index": 23
  },
  {
    "timestamp": 1546623267000,
    "block_number": 7010037,
    "tx_hash": "0xcb91ef8114bdeafafb12b47ca9f9b76d2cc1db9dba655cb8a4aa339afc65abb5",
    "eth_amount": 60528786878456780,
    "user_addr": "0x41cc8e0abcd7d4f2d06464a2c69789064dee42f5",
    "src_addr": "0x094c875704c14783049ddf8136e298b3a099c446",
    "dst_addr": "0xdd974d5c2e2928dea5f71b9825b8b646686bd200",
    "src_reserve_addr": "0x2295fc6bc32cd12fdbb852cff4014ceac6d79c10",
    "dst_reserve_addr": "0x0000000000000000000000000000000000000000",
    "src_amount": 8885400000000001000,
    "dst_amount": 57520001928530895000,
    "fiat_amount": 8.881425181696164,
    "wallet_addr": "0xea1a7de54a427342c8820185867cf49fc2f95d43",
    "src_burn_amount": 0.1006291081854344,
    "dst_burn_amount": 0,
    "src_wallet_fee_amount": 0,
    "dst_wallet_fee_amount": 0,
    "integration_app": "KyberSwap",
    "ip": "1.246.178.89",
    "country": "KR",
    "user_name": "",
    "profile_id": 0,
    "index": 43
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
